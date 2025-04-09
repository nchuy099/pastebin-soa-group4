package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"get-paste-service/model"
	"get-paste-service/repository"

	"github.com/julienschmidt/httprouter"
)

func GetPasteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pasteID := ps.ByName("id")

	paste, err := repository.GetPasteByID(pasteID)
	if err != nil {
		if err == repository.ErrPasteExpired || err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Paste not found or expired")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
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
func respondWithError(w http.ResponseWriter, code int, message string) {
	response := model.ResponseData{
		Status:  code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
