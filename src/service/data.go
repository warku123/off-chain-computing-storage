package service

import (
	"encoding/json"
	"io"
	"net/http"
	"offstorage/data"
	"os"
	"path"
)

type DataSession struct {
	Dsh *data.Data_api // shell instance
}

var data_path string = "/Users/jojo/test/data"

func CreateData(w http.ResponseWriter, r *http.Request) {
	if !pathExists(data_path) {
		os.MkdirAll(data_path, 0755)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Create a map to hold the parsed JSON data
	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	task_id, ok := requestData["taskID"].(string)
	if !ok {
		task_id = ""
	}

	local_path, ok := requestData["localPath"].(string)
	if !ok {
		local_path = "/data"
	}

	dsh, task_id, err := data.NewDataShell(
		data.DataWithHost("127.0.0.1"),
		data.DataWithPort(5001),
		data.DataWithKeyName(requestData["keyName"].(string)),
		data.DataWithIpnsName(requestData["ipnsName"].(string)),
		data.DataWithLocalPath(local_path),
		data.DataWithRole(requestData["role"].(string)),
		data.DataWithTaskID(task_id),
	)
	if err != nil {
		http.Error(w, "Error creating data shell: "+err.Error(), http.StatusBadRequest)
		return
	}

	sessionID, err := GenerateRandomID(16)
	if err != nil {
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	newSession := &DataSession{
		Dsh: dsh,
	}

	sessionStore.Store(sessionID, newSession)

	responseData := map[string]string{
		"session_id": sessionID,
		"task_id":    task_id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func CatData(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID := r.URL.Query().Get("session_id")

	// Get the session from the session store using sync.Map's Load method (thread-safe)
	sessionObj, ok := sessionStore.Load(sessionID)
	if !ok {
		http.Error(w, "Error getting session", http.StatusBadRequest)
		return
	}

	// Type assertion to convert the interface{} type to *Session
	sessionData := sessionObj.(*DataSession)

	// Parse incoming JSON data
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Get the data
	data, err := sessionData.Dsh.CatData(requestData["key"].(string))
	if err != nil {
		http.Error(w, "Error getting data", http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"data": data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*DataSession)
	if !ok {
		http.Error(w, "Invalid session object", http.StatusInternalServerError)
		return
	}

	// Parse incoming JSON data
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Get the data from the request
	name := requestData["name"].(string)
	destFilePath := path.Join(data_path, name)
	err = session.Dsh.GetData(name, destFilePath)
	if err != nil {
		http.Error(w, "Failed to get data:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(destFilePath)

	// 设置Content-Disposition头，指示浏览器下载文件而不是直接在浏览器中打开
	w.Header().Set("Content-Disposition", "attachment; filename="+name)

	// 使用http.ServeFile函数将文件内容作为响应发送给客户端
	http.ServeFile(w, r, destFilePath)
}

func AddDataString(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID := r.URL.Query().Get("session_id")

	// Get the session from the session store using sync.Map's Load method (thread-safe)
	sessionObj, ok := sessionStore.Load(sessionID)
	if !ok {
		http.Error(w, "Error getting session", http.StatusBadRequest)
		return
	}

	// Type assertion to convert the interface{} type to *Session
	sessionData := sessionObj.(*DataSession)

	// Parse incoming JSON data
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Add the data
	cid, err := sessionData.Dsh.AddDataString(
		requestData["key"].(string),
		requestData["value"].(string),
	)
	if err != nil {
		http.Error(w, "Error adding data", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	responseData := map[string]string{
		"cid": cid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func AddDataFile(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*DataSession)
	if !ok {
		http.Error(w, "Invalid session object", http.StatusInternalServerError)
		return
	}

	// 解析文件上传
	err := r.ParseMultipartForm(10 << 20) // 最大支持10MB的文件上传
	if err != nil {
		http.Error(w, "Failed to parse multipart form data", http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "Key need to be provided in form value", http.StatusBadRequest)
		return
	}

	// 保存文件到本地
	destFilePath := path.Join(data_path, key)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		http.Error(w, "Failed to create destination file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(destFilePath)
	defer destFile.Close()

	// 将文件内容拷贝到目标文件
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Failed to copy file content", http.StatusInternalServerError)
		return
	}

	cid, err := session.Dsh.AddDataFile(key, destFilePath)
	if err != nil {
		http.Error(w, "Error adding data", http.StatusInternalServerError)
		return
	}

	// Respond
	responseData := map[string]string{
		"cid": cid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func DataPersistant(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID := r.URL.Query().Get("session_id")

	// Get the session from the session store using sync.Map's Load method (thread-safe)
	sessionObj, ok := sessionStore.Load(sessionID)
	if !ok {
		http.Error(w, "Error getting session", http.StatusBadRequest)
		return
	}

	// Type assertion to convert the interface{} type to *Session
	session := sessionObj.(*DataSession)

	err := session.Dsh.DataPersistance()
	if err != nil {
		http.Error(w, "Error persisting data", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	responseData := map[string]string{
		"message": "Data persisted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func CloseData(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID := r.URL.Query().Get("session_id")

	// Get the session from the session store using sync.Map's Load method (thread-safe)
	sessionObj, ok := sessionStore.Load(sessionID)
	if !ok {
		http.Error(w, "Error getting session", http.StatusBadRequest)
		return
	}

	// Type assertion to convert the interface{} type to *Session
	session := sessionObj.(*DataSession)

	// Close the session
	if session.Dsh != nil {
		role := session.Dsh.GetRole()
		var err error
		if role == "judge" {
			err = session.Dsh.CloseJudgeData()
		} else {
			err = session.Dsh.CloseData()
		}
		if err != nil {
			http.Error(w, "Error closing data shell", http.StatusInternalServerError)
			return
		}
	}

	sessionStore.Delete(sessionID)

	// Respond with a success message
	responseData := map[string]string{
		"message": "Session closed and removed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
