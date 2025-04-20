package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"get-stats-service/cache"
	"get-stats-service/model"
	"get-stats-service/repository"

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

	// Get current year and month
	now := time.Now()
	currentYear, currentMonth := now.Year(), int(now.Month())

	if year > currentYear || (year == currentYear && month > currentMonth) {
		respondWithError(w, http.StatusBadRequest, "Requested year and month cannot be in the future", errors.New("invalid year and month"))
		log.Printf("Requested year and month cannot be in the future: %d-%02d", year, month)
		return
	}

	// Check if requested stats are not for the current month
	if year != currentYear || month != currentMonth {
		// Try to get from cache first for non-current months
		cachedStats, err := cache.GetMonthlyStatsFromCache(year, month)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch stats", err)
			return
		} else if cachedStats != nil {
			response := model.ResponseData{
				Status:  http.StatusOK,
				Message: "Get monthly stats successfully",
				Data:    cachedStats,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}

	// Get stats from repository
	var stats *model.MonthlyStats
	if year == currentYear && month == currentMonth {
		stats, err = repository.GetMonthlyStatsForCurrentMonth(year, month)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch stats", err)
			return
		}
	} else {
		stats, err = repository.GetMonthlyStatsForNonCurrentMonth(year, month)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch stats", err)
			return
		}
	}

	// If not current month, cache the results
	if year != currentYear || month != currentMonth {
		if err := cache.SetMonthlyStatsToCache(year, month, stats); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch stats", err)
			return
		}
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
