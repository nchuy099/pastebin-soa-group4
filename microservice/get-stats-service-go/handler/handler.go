package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"get-stats-service-go/model"
	"get-stats-service-go/repository"

	"github.com/julienschmidt/httprouter"
)

// GetMonthlyStats handles requests for monthly paste statistics
func GetMonthlyStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the year and month from query parameters
	yearMonth := r.URL.Query().Get("month")
	if yearMonth == "" {
		respondWithError(w, http.StatusBadRequest, "Month parameter is required (format: YYYY-MM)", errors.New("invalid month format"))
		log.Println("Month parameter is required (format: YYYY-MM)")
		return
	}

	// Parse the year-month parameter
	date, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid month format. Expected YYYY-MM", err)
		log.Printf("Invalid month format. Expected YYYY-MM, got: %s", yearMonth)
		return
	}

	year, month := date.Year(), int(date.Month())

	// Get stats from repository
	stats, err := repository.GetMonthlyStats(year, month)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch stats", err)
		return
	}

	// Create response
	response := model.ResponseData{
		Status:  http.StatusOK,
		Message: "Get monthly stats successfully",
		Data:    stats,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	ersMsg := err.Error()
	response := model.ResponseData{
		Status:  code,
		Message: message,
		Error:   &ersMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
