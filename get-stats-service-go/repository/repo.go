package repository

import (
	"get-stats-service-go/db"
	"get-stats-service-go/model"
	"log"
	"time"
)

// GetMonthlyStats retrieves monthly statistics for pastes created in a specific year and month
func GetMonthlyStats(year, month int) (*model.MonthlyStats, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // Tháng kế tiếp

	query := `
		SELECT
			COALESCE(SUM(views), 0) AS total_views,
			COALESCE(AVG(views), 0) AS avg_views_per_paste,
			COALESCE(MIN(views), 0) AS min_views,
			COALESCE(MAX(views), 0) AS max_views
		FROM paste
		WHERE created_at >= ? AND created_at < ?
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
