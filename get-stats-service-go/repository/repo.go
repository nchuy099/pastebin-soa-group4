package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"get-stats-service/db"
	"get-stats-service/model"
	"log"
	"time"
)

// GetMonthlyStatsForNonCurrentMonth fetches stats directly from the materialized view database
func GetMonthlyStatsForNonCurrentMonth(year, month int) (*model.MonthlyStats, error) {
	monthStr := fmt.Sprintf("%04d-%02d", year, month) // e.g. "2025-04"

	query := `
        SELECT
            total_views,
            avg_views_per_paste,
            min_views,
            max_views
        FROM paste_stats_readonly
        WHERE month_year = ?
    `

	row := db.DB.QueryRow(query, monthStr)

	var stats model.MonthlyStats
	err := row.Scan(
		&stats.TotalViews,
		&stats.AvgViewsPerPaste,
		&stats.MinViews,
		&stats.MaxViews,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Không có dữ liệu thống kê cho tháng này
			return &model.MonthlyStats{
				TotalViews:       0,
				AvgViewsPerPaste: 0,
				MinViews:         0,
				MaxViews:         0,
			}, nil
		}

		log.Printf("Error scanning cached monthly stats: %v", err.Error())
		return nil, err
	}

	return &stats, nil
}

// GetMonthlyStatsForCurrentMonth fetches stats directly from the primary database
func GetMonthlyStatsForCurrentMonth(year, month int) (*model.MonthlyStats, error) {
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
