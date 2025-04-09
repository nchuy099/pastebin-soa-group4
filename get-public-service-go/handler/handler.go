package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"get-public-service/model"
	"get-public-service/repository"

	"github.com/julienschmidt/httprouter"
)

func GetPublicPastes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Lấy limit và page từ query params
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Default page
	}

	offset := (page - 1) * limit // Tính offset

	// Fetch pastes with pagination
	pastes, err := repository.GetPublicPastes(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch pastes")
		return
	}

	// Get total count for pagination
	totalItems, err := repository.CountPublicPastes()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to count pastes")
		return
	}

	// Create response with pagination
	pasteListResponse := model.PasteListResponse{
		Pastes: pastes,
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
		},
	}

	response := model.ResponseData{
		Status:  http.StatusOK,
		Message: "Public pastes retrieved successfully",
		Data:    pasteListResponse,
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
