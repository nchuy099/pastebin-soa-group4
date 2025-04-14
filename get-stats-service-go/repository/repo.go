package repository

import (
	"log"

	"get-stats-service-go/db"
	"get-stats-service-go/model"
)

// GetMonthlyStats retrieves monthly statistics for pastes created in a specific year and month
func GetMonthlyStats(year, month int) (*model.MonthlyStats, error) {
	query := `
		SELECT
			COUNT(id) AS total_pastes,
			COALESCE(SUM(views), 0) AS total_views,
			COALESCE(AVG(views), 0) AS avg_views_per_paste,
			COALESCE(MIN(views), 0) AS min_views,
			COALESCE(MAX(views), 0) AS max_views,
			COALESCE(SUM(CASE WHEN expires_at IS NULL OR expires_at > NOW() THEN 1 ELSE 0 END), 0) AS active_pastes,
			COALESCE(SUM(CASE WHEN expires_at IS NOT NULL AND expires_at <= NOW() THEN 1 ELSE 0 END), 0) AS expired_pastes
		FROM paste
		WHERE YEAR(created_at) = ? AND MONTH(created_at) = ?
	`

	row := db.DB.QueryRow(query, year, month)

	var stats model.MonthlyStats
	err := row.Scan(
		&stats.TotalPastes,
		&stats.TotalViews,
		&stats.AvgViewsPerPaste,
		&stats.MinViews,
		&stats.MaxViews,
		&stats.ActivePastes,
		&stats.ExpiredPastes,
	)

	if err != nil {
		log.Printf("Error scanning monthly stats: %v", err.Error())
		return nil, err
	}

	return &stats, nil
}
