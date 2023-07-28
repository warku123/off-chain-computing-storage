package service

import (
	"encoding/json"
	"io"
	"net/http"
	"offstorage/data"
)

type DataSession struct {
	Dsh *data.Data_api // shell instance
}

func CreateData(w http.ResponseWriter, r *http.Request) {
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

func GetData(w http.ResponseWriter, r *http.Request) {
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
	data, err := sessionData.Dsh.GetData(requestData["key"].(string))
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

func AddData(w http.ResponseWriter, r *http.Request) {
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
	err = sessionData.Dsh.AddData(
		requestData["key"].(string),
		requestData["value"].(string),
	)
	if err != nil {
		http.Error(w, "Error adding data", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	responseData := map[string]string{
		"message": "Data added successfully",
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
