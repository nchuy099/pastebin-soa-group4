package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"create-paste-service-go/model"
	"create-paste-service-go/repository"

	"github.com/julienschmidt/httprouter"
)

// GenerateUniqueID generates a unique ID for a paste
func GenerateUniqueID(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating ID: %v", err.Error())
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreatePaste handles the creation of a new paste
func CreatePaste(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse request body
	var request model.CreatePasteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		log.Printf("Error decoding request body: %v", err.Error())
		return
	}
	defer r.Body.Close()

	// Generate unique ID
	var id string
	var err error
	var exists bool

	for {
		id, err = GenerateUniqueID(8)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate ID", err)
			return
		}

		exists, err = repository.ExistsById(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to check ID availability", err)
			return
		}

		if !exists {
			break
		}
	}

	// Set defaults if values are empty
	if request.Title == "" {
		request.Title = "Untitled"
	}

	if request.Language == "" {
		request.Language = "text"
	}

	// Set visibility
	visibility := model.Public
	if strings.ToUpper(request.Visibility) == string(model.Unlisted) {
		visibility = model.Unlisted
	}

	// Calculate expiresAt
	var expiresAt *time.Time
	if request.ExpiresIn > 0 {
		expires := time.Now().Add(time.Duration(request.ExpiresIn) * time.Minute)
		expiresAt = &expires
	}

	// Create paste object
	paste := &model.Paste{
		ID:         id,
		Content:    request.Content,
		Title:      request.Title,
		Language:   request.Language,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		Views:      0,
		Visibility: visibility,
	}

	// Save to database
	if err := repository.SavePaste(paste); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save paste", err)
		return
	}

	// Return success response
	response := model.ResponseData{
		Status:  http.StatusOK,
		Message: "Paste created successfully",
		Data: struct {
			ID string `json:"id"`
		}{ID: id},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	errMsg := err.Error()
	response := model.ResponseData{
		Status:  code,
		Message: message,
		Error:   &errMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
