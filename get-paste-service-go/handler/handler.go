package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"get-paste-service/cache"
	"get-paste-service/model"
	"get-paste-service/producer"
	"get-paste-service/repository"

	"github.com/julienschmidt/httprouter"
)

func GetPasteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pasteID := ps.ByName("id")

	paste := &model.Paste{}
	var expiresAt sql.NullTime

	// Try to get paste from cache first
	paste, err := cache.GetPasteFromCache(pasteID)

	// If not in cache, get from database
	log.Printf("Cache miss for paste ID: %s, fetching from database", pasteID)

	paste, err = repository.GetPasteByID(pasteID)
	if err != nil {
		if err == repository.ErrPasteExpired {
			respondWithError(w, http.StatusForbidden, "Paste expired", err)
		} else if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Paste not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		}
		return
	}

	if expiresAt.Valid {
		paste.ExpiresAt = &expiresAt.Time
	}

	// Store paste in cache for future requests
	err = cache.SetPasteToCache(paste)
	if err != nil {
		log.Printf("Failed to cache paste: %v", err)
	}

	err = producer.PublishPasteView(&model.PasteView{
		PasteID:  pasteID,
		ViewedAt: time.Now().UTC(),
	})

	if err != nil {
		// Just log the error but continue to serve the paste
		// This way view counts might be slightly off but the user experience is not affected
		// The worker will retry publishing view updates
		log.Printf("Error recording paste view: %v", err)
	}

	response := model.ResponseData{
		Status:  http.StatusOK,
		Message: "Paste retrieved successfully",
		Data:    paste,
	}

	w.Header().Set("Content-Type", "application/json")
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
