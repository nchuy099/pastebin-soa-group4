package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"get-stats-service/model"
	"get-stats-service/repository"

	"github.com/julienschmidt/httprouter"
)

// Định nghĩa múi giờ GMT+7 (giống với repository package)
var gmt7 = time.FixedZone("GMT+7", 7*60*60)

// GetStats handles requests for paste view statistics
func GetStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the mode from query parameters
	mode := r.URL.Query().Get("mode")

	var stats *model.Stats
	var err error

	// Convert to GMT+7 timezone
	currentTime := time.Now().In(gmt7)

	switch mode {
	case "last-10-minutes":
		stats, err = repository.GetStatsForLast10Minutes(currentTime)
	case "last-24-hours":
		stats, err = repository.GetStatsForLastDay(currentTime)
	case "last-7-days":
		stats, err = repository.GetStatsForLastWeek(currentTime)
	case "last-30-days":
		stats, err = repository.GetStatsForLastMonth(currentTime)
	default:
		// Default to last 10 minutes if no mode specified
		stats, err = repository.GetStatsForLast10Minutes(currentTime)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get stats", err)
		log.Println("Error getting stats:", err)
		return
	}

	// Create response
	response := model.ResponseData{
		Status:  http.StatusOK,
		Message: "Stats retrieved successfully (GMT+7 timezone)",
		Data: model.Stats{
			TimeViews:  stats.TimeViews,
			TotalViews: stats.TotalViews,
			Timezone:   "GMT+7",
		},
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
