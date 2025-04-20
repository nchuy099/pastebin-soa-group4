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
		if err == repository.ErrPasteExpired {
			respondWithError(w, http.StatusForbidden, "Paste expired", err)
		} else if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Paste not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error", err)
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
