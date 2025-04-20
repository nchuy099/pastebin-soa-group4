package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"get-paste-service/model"
	"get-paste-service/repository"

	"github.com/julienschmidt/httprouter"
)

func GetPasteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pasteID := ps.ByName("id")

	paste, err := repository.GetPasteByID(pasteID)
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

	// Add paste view asynchronously
	remoteAddr := r.RemoteAddr
	userAgent := r.UserAgent()

	err = repository.AddPasteView(&model.PasteViews{
		PasteID:  pasteID,
		ViewedAt: time.Now().UTC(),
	}, remoteAddr, userAgent)

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
