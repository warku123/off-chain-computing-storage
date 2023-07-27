package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"offstorage/image"
	"os"
	"time"
)

// Custom data structure to represent the session
type ImageSession struct {
	Ish *image.Image_api // shell instance
}

func CreateImage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
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

	ish, err := image.NewImageShell(
		image.ImageWithHost("127.0.0.1"),
		image.ImageWithPort(5001),
		image.ImageWithKeyName(data["keyName"].(string)),
		image.ImageWithIpnsName(data["ipnsName"].(string)),
		image.ImageWithLocalPath("/image"),
	)

	if err != nil {
		http.Error(w, "Error creating image shell", http.StatusInternalServerError)
		return
	}

	sessionID, err := GenerateRandomID(16)
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

func GetImageByCid(w http.ResponseWriter, r *http.Request) {
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

func GetImageByIdx(w http.ResponseWriter, r *http.Request) {
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
	idx, ok := requestData["idx"].(int)
	if !ok {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	imageData, timestamp, err := sessionData.Ish.CatImageByIdx(idx)
	if err != nil {
		http.Error(w, "Error getting image", http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"image_data": imageData,
		"timestamp":  timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func AddImage(w http.ResponseWriter, r *http.Request) {
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
	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Get the data from the request
	data := requestData["data"]

	tmp_path := "/tmp_image.txt"
	err = os.WriteFile(tmp_path, []byte(data), 0644)
	if err != nil {
		http.Error(w, "Failed to write data to IPFS", http.StatusInternalServerError)
		return
	}

	timestamp := time.Now().Unix()

	// Use IPFS Shell's Add method to write data to IPFS
	cid, idx, err := session.Ish.AddImage(tmp_path, fmt.Sprint(timestamp))
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

// CloseSession closes a session and removes it from the sessionStore
func CloseImageSession(w http.ResponseWriter, r *http.Request) {
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
