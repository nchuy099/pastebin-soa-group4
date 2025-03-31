package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"get-paste-service/repository"
)

func GetPasteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pasteID := ps.ByName("id")

	paste, err := repository.GetPasteByID(pasteID)
	if err != nil {
		if err == repository.ErrPasteExpired || err == sql.ErrNoRows {
			http.Error(w, `{"error": "Paste not found or expired"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paste)
}
