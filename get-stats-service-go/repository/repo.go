package repository

import (
	"get-stats-service-go/cache"
	"get-stats-service-go/db"
	"get-stats-service-go/model"
	"log"
	"time"
)

// GetMonthlyStats retrieves monthly statistics for pastes created in a specific year and month
func GetMonthlyStats(year, month int) (*model.MonthlyStats, error) {
	// Get current year and month
	now := time.Now()
	currentYear, currentMonth := now.Year(), int(now.Month())

	// Check if requested stats are not for the current month
	if year != currentYear || month != currentMonth {
		// Try to get from cache first for non-current months
		log.Printf("Attempting to fetch stats for %d-%02d from cache (not current month)", year, month)
		cachedStats, err := cache.GetMonthlyStatsFromCache(year, month)
		if err != nil {
			log.Printf("Error getting stats from cache: %v", err)
		} else if cachedStats != nil {
			return cachedStats, nil
		}
	}

	// Continue with database query if not in cache or is current month
	log.Printf("Fetching stats for %d-%02d from database", year, month)
	stats, err := getMonthlyStatsFromDB(year, month)
	if err != nil {
		return nil, err
	}

	// If not current month, cache the results
	if year != currentYear || month != currentMonth {
		if err := cache.SetMonthlyStatsToCache(year, month, stats); err != nil {
			log.Printf("Failed to cache stats: %v", err)
		}
	} else {
		// For current month, still cache but with short TTL (handled in SetMonthlyStatsToCache)
		if err := cache.SetMonthlyStatsToCache(year, month, stats); err != nil {
			log.Printf("Failed to cache current month stats: %v", err)
		}
	}

	return stats, nil
}

// getMonthlyStatsFromDB fetches stats directly from the database
func getMonthlyStatsFromDB(year, month int) (*model.MonthlyStats, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // Next month

	query := `
            SELECT
                COALESCE(SUM(view_count), 0) AS totalViews,
                COALESCE(ROUND(AVG(view_count)), 0) AS avgViewsPerPaste,
                COALESCE(MIN(view_count), 0) AS minViews,
                COALESCE(MAX(view_count), 0) AS maxViews
            FROM (
                SELECT paste_id, COUNT(*) AS view_count
                FROM paste_views
                WHERE viewed_at >= ? AND viewed_at < ?
                GROUP BY paste_id
            ) AS monthly_paste_views
	`

	row := db.DB.QueryRow(query, startDate, endDate)

	var stats model.MonthlyStats
	err := row.Scan(
		&stats.TotalViews,
		&stats.AvgViewsPerPaste,
		&stats.MinViews,
		&stats.MaxViews,
	)

	if err != nil {
		log.Printf("Error scanning monthly stats: %v", err.Error())
		return nil, err
	}

	return &stats, nil
}
