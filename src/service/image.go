package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"offstorage/image"
	"offstorage/utils"
	"os"
	"path"
	"strconv"
)

// Custom data structure to represent the session
type ImageSession struct {
	Ish *image.Image_api // shell instance
}

var image_path string = "/Users/jojo/test/image"

func CreateImage(w http.ResponseWriter, r *http.Request) {
	if !utils.PathExists(image_path) {
		os.MkdirAll(image_path, 0755)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Create a map to hold the parsed JSON data
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	local_path, ok := data["localPath"].(string)
	if !ok {
		local_path = "/image"
	}

	ish, err := image.NewImageShell(
		image.ImageWithHost("127.0.0.1"),
		image.ImageWithPort(5001),
		image.ImageWithKeyName(data["keyName"].(string)),
		image.ImageWithIpnsName(data["ipnsName"].(string)),
		image.ImageWithTaskName(data["taskName"].(string)),
		image.ImageWithOwnerName(data["ownerName"].(string)),
		image.ImageWithLocalPath(local_path),
	)

	if err != nil {
		http.Error(w, "Error creating image shell: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := utils.GenerateRandomID(16)
	if err != nil {
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	newSession := &ImageSession{
		Ish: ish,
	}

	// Save the session in the session store using sync.Map's Store method (thread-safe)
	sessionStore.Store(sessionID, newSession)

	responseData := map[string]string{
		"session_id": sessionID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func CatImageByCid(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID := r.URL.Query().Get("session_id")

	// Get the session from the session store using sync.Map's Load method (thread-safe)
	sessionObj, ok := sessionStore.Load(sessionID)
	if !ok {
		http.Error(w, "Error getting session", http.StatusBadRequest)
		return
	}

	// Type assertion to convert the interface{} type to *Session
	sessionData := sessionObj.(*ImageSession)

	// Parse incoming JSON data
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Get the IPFS hash from the request data
	cid, ok := requestData["cid"].(string)
	if !ok {
		http.Error(w, "Invalid CID", http.StatusBadRequest)
		return
	}

	imageData, err := sessionData.Ish.CatImageByCid(cid)
	if err != nil {
		http.Error(w, "Error getting image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imageData)
}

func CatImageByIdx(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID := r.URL.Query().Get("session_id")

	// Get the session from the session store using sync.Map's Load method (thread-safe)
	sessionObj, ok := sessionStore.Load(sessionID)
	if !ok {
		http.Error(w, "Error getting session", http.StatusBadRequest)
		return
	}

	// Type assertion to convert the interface{} type to *Session
	sessionData := sessionObj.(*ImageSession)

	// Parse incoming JSON data
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// 数字在interface里面必须先转为float64，再变为int
	idx, ok := requestData["idx"].(float64)
	if !ok {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	imageData, height, err := sessionData.Ish.CatImageByIdx(int(idx))
	if err != nil {
		http.Error(w, "Error getting image", http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"image_data": imageData,
		"height":     height,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func GetImageByCid(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
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
	cid := requestData["cid"].(string)
	destFilePath := path.Join(image_path, cid)
	err = session.Ish.GetImageByCid(cid, destFilePath)
	if err != nil {
		http.Error(w, "Failed to get image:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(destFilePath)

	// 设置Content-Disposition头，指示浏览器下载文件而不是直接在浏览器中打开
	w.Header().Set("Content-Disposition", "attachment; filename="+cid)

	// 使用http.ServeFile函数将文件内容作为响应发送给客户端
	http.ServeFile(w, r, destFilePath)
}

func GetImageByIdx(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
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
	idx, ok := requestData["idx"].(float64)
	if !ok {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	destFilePath := path.Join(image_path, fmt.Sprint(int(idx)))
	_, err = session.Ish.GetImageByIdx(int(idx), destFilePath)
	if err != nil {
		http.Error(w, "Failed to get image:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(destFilePath)

	// 设置Content-Disposition头，指示浏览器下载文件而不是直接在浏览器中打开
	w.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprint(int(idx)))

	// 使用http.ServeFile函数将文件内容作为响应发送给客户端
	http.ServeFile(w, r, destFilePath)
}

func GetHeightByIdx(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
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
	idx, ok := requestData["idx"].(float64)
	if !ok {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	height, err := session.Ish.GetHeight(int(idx))
	if err != nil {
		http.Error(w, "Failed to get height:"+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"height": height,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// AddImageString adds image to IPFS，data是base64编码后的字符串
func AddImageString(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
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
	data := requestData["data"].(string)

	tmp_path := path.Join(image_path, "tmp")
	defer os.Remove(tmp_path)

	err = os.WriteFile(tmp_path, []byte(data), 0644)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	height, ok := requestData["height"].(float64)
	if !ok {
		http.Error(w, "Invalid height", http.StatusBadRequest)
		return
	}

	// Use IPFS Shell's Add method to write data to IPFS
	cid, idx, err := session.Ish.AddImage(tmp_path, uint64(height))
	if err != nil {
		http.Error(w, "Failed to write data to IPFS", http.StatusInternalServerError)
		return
	}

	// Respond with the IPFS hash
	responseData := map[string]interface{}{
		"cid": cid,
		"idx": idx,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func AddImageFile(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
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

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 保存文件到本地
	destFilePath := path.Join(image_path, fileHeader.Filename)
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

	s_height := r.Form.Get("height")
	height, err := strconv.ParseUint(s_height, 10, 64)
	if err != nil {
		http.Error(w, "Invalid height: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Use IPFS Shell's Add method to write data to IPFS
	cid, idx, err := session.Ish.AddImage(destFilePath, height)
	if err != nil {
		http.Error(w, "Failed to write data to IPFS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the IPFS hash
	responseData := map[string]interface{}{
		"cid": cid,
		"idx": idx,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func GetMTreeRootHash(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
	if !ok {
		http.Error(w, "Invalid session object", http.StatusInternalServerError)
		return
	}

	hash, err := session.Ish.GetRootHash()
	if err != nil {
		http.Error(w, "Error getting root hash", http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"hash": hash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// CloseSession closes a session and removes it from the sessionStore
func CloseImage(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	// Load the session from the sessionStore
	sessionObj, exists := sessionStore.Load(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Type assertion to retrieve the session object from sync.Map
	session, ok := sessionObj.(*ImageSession)
	if !ok {
		http.Error(w, "Invalid session object", http.StatusInternalServerError)
		return
	}

	// Close the IPFS Shell instance
	if session.Ish != nil {
		err := session.Ish.CloseImage()
		if err != nil {
			http.Error(w, "Error closing image shell", http.StatusInternalServerError)
			return
		}
	}

	// Remove the session from the sessionStore
	sessionStore.Delete(sessionID)

	// Respond with a success message
	responseData := map[string]string{
		"message": "Session closed and removed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
