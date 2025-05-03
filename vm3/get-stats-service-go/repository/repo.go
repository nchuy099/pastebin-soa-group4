package repository

import (
	"get-stats-service/db"
	"get-stats-service/model"
	"log"
	"time"
)

// Định nghĩa múi giờ GMT+7
var gmt7 = time.FixedZone("GMT+7", 7*60*60)

func GetStatsForLast10Minutes(pasteID string, currentTime time.Time) (*model.Stats, error) {
	currentTime = currentTime.In(gmt7)

	// Không làm tròn endTime để bao gồm phút hiện tại
	endTime := currentTime
	startTime := endTime.Add(-10 * time.Minute)

	query := `
		SELECT 
			DATE_FORMAT(CONVERT_TZ(viewed_at, '+00:00', '+07:00'), '%Y-%m-%d %H:%i') AS time_label,
			COUNT(*) AS views
		FROM paste_views
		WHERE paste_id = ? AND viewed_at >= ? AND viewed_at < ?
		GROUP BY time_label
		ORDER BY time_label ASC
	`

	rows, err := db.DB.Query(query, pasteID, startTime.UTC(), endTime.UTC())
	if err != nil {
		log.Printf("Error querying DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Khởi tạo map để ánh xạ thời gian -> số lượt xem
	viewMap := make(map[string]int64)

	var totalViews int64
	for rows.Next() {
		var timeLabel string
		var views int64
		if err := rows.Scan(&timeLabel, &views); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		viewMap[timeLabel] = views // Sử dụng timeLabel làm khóa
		totalViews += views
	}

	// Tạo danh sách đầy đủ cho 10 phút, bắt đầu từ phút hiện tại
	var timeViews []model.TimeView
	for i := 9; i >= 0; i-- {
		t := endTime.Add(-time.Duration(i) * time.Minute).Truncate(time.Minute)
		timeKey := t.Format("2006-01-02 15:04") // Định dạng khớp với time_label
		label := t.Format("15:04")              // Định dạng hiển thị, ví dụ: "19:24"
		views := viewMap[timeKey]               // Lấy số lượt xem từ viewMap, mặc định 0 nếu không có
		timeViews = append(timeViews, model.TimeView{
			Time:  label,
			Views: views,
		})
	}

	return &model.Stats{
		PasteID:    pasteID,
		TimeViews:  timeViews,
		TotalViews: totalViews,
	}, nil
}

func GetStatsForLastDay(pasteID string, currentTime time.Time) (*model.Stats, error) {
	// Chuyển đổi thời gian hiện tại sang GMT+7
	currentTime = currentTime.In(gmt7)

	endTime := currentTime.Truncate(time.Hour)
	startTime := endTime.Add(-24 * time.Hour)

	query := `
		SELECT 
			DATE_FORMAT(CONVERT_TZ(viewed_at, '+00:00', '+07:00'), '%Y-%m-%d %H:00') AS time_label,
			COUNT(*) AS views
		FROM paste_views
		WHERE paste_id = ? AND viewed_at >= ? AND viewed_at < ?
		GROUP BY time_label
		ORDER BY time_label ASC
	`

	rows, err := db.DB.Query(query, pasteID, startTime.UTC(), endTime.UTC())
	if err != nil {
		log.Printf("Error querying DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Khởi tạo map để ánh xạ thời gian -> số lượt xem
	viewMap := make(map[string]int64)

	var totalViews int64
	for rows.Next() {
		var label string
		var views int64
		if err := rows.Scan(&label, &views); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		viewMap[label] = views
		totalViews += views
	}

	// Tạo danh sách đầy đủ cho 24 giờ
	var timeViews []model.TimeView
	for i := 23; i >= 0; i-- {
		t := endTime.Add(-time.Duration(i) * time.Hour)
		label := t.Format("15:04")
		hourKey := t.Format("2006-01-02 15:00")
		views := viewMap[hourKey]
		timeViews = append(timeViews, model.TimeView{
			Time:  label,
			Views: views,
		})
	}

	return &model.Stats{
		PasteID:    pasteID,
		TimeViews:  timeViews,
		TotalViews: totalViews,
	}, nil
}

func GetStatsForLastWeek(pasteID string, currentTime time.Time) (*model.Stats, error) {
	// Chuyển đổi thời gian hiện tại sang GMT+7
	currentTime = currentTime.In(gmt7)

	endTime := currentTime.Truncate(24 * time.Hour)
	startTime := endTime.AddDate(0, 0, -7)

	query := `
		SELECT 
			DATE_FORMAT(CONVERT_TZ(viewed_at, '+00:00', '+07:00'), '%Y-%m-%d') AS time_label,
			COUNT(*) AS views
		FROM paste_views
		WHERE paste_id = ? AND viewed_at >= ? AND viewed_at < ?
		GROUP BY time_label
		ORDER BY time_label ASC
	`

	rows, err := db.DB.Query(query, pasteID, startTime.UTC(), endTime.UTC())
	if err != nil {
		log.Printf("Error querying DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Khởi tạo map để ánh xạ thời gian -> số lượt xem
	viewMap := make(map[string]int64)

	var totalViews int64
	for rows.Next() {
		var label string
		var views int64
		if err := rows.Scan(&label, &views); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		viewMap[label] = views
		totalViews += views
	}

	// Tạo danh sách đầy đủ cho 7 ngày
	var timeViews []model.TimeView
	for i := 6; i >= 0; i-- {
		t := endTime.AddDate(0, 0, -i)
		label := t.Format("01/02")
		dayKey := t.Format("2006-01-02")
		views := viewMap[dayKey]
		timeViews = append(timeViews, model.TimeView{
			Time:  label,
			Views: views,
		})
	}

	return &model.Stats{
		PasteID:    pasteID,
		TimeViews:  timeViews,
		TotalViews: totalViews,
	}, nil
}

func GetStatsForLastMonth(pasteID string, currentTime time.Time) (*model.Stats, error) {
	// Chuyển đổi thời gian hiện tại sang GMT+7
	currentTime = currentTime.In(gmt7)

	endTime := currentTime.Truncate(24 * time.Hour)
	startTime := endTime.AddDate(0, 0, -30)

	query := `
		SELECT 
			DATE_FORMAT(CONVERT_TZ(viewed_at, '+00:00', '+07:00'), '%Y-%m-%d') AS time_label,
			COUNT(*) AS views
		FROM paste_views
		WHERE paste_id = ? AND viewed_at >= ? AND viewed_at < ?
		GROUP BY time_label
		ORDER BY time_label ASC
	`

	rows, err := db.DB.Query(query, pasteID, startTime.UTC(), endTime.UTC())
	if err != nil {
		log.Printf("Error querying DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Khởi tạo map để ánh xạ thời gian -> số lượt xem
	viewMap := make(map[string]int64)

	var totalViews int64
	for rows.Next() {
		var label string
		var views int64
		if err := rows.Scan(&label, &views); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		viewMap[label] = views
		totalViews += views
	}

	// Tạo danh sách đầy đủ cho 30 ngày
	var timeViews []model.TimeView
	for i := 29; i >= 0; i-- {
		t := endTime.AddDate(0, 0, -i)
		label := t.Format("01/02") // Hiển thị theo định dạng MM/DD
		dayKey := t.Format("2006-01-02")
		views := viewMap[dayKey]
		timeViews = append(timeViews, model.TimeView{
			Time:  label,
			Views: views,
		})
	}

	return &model.Stats{
		PasteID:    pasteID,
		TimeViews:  timeViews,
		TotalViews: totalViews,
	}, nil
}

func GetTotalPasteViewsFromCreation(pasteID string) (int64, error) {
	query := `
		SELECT COUNT(*) AS total_views
		FROM paste_views
		WHERE paste_id = ?;	
	`

	var totalViews int64
	err := db.DB.QueryRow(query, pasteID).Scan(&totalViews)
	if err != nil {
		return 0, err
	}

	return totalViews, nil
}
