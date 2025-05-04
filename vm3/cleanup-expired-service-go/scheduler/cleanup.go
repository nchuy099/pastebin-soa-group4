package scheduler

import (
	"cleanup-expired-service-go/db"
	"log"
	"time"
)

// CleanupExpiredPastes deletes all pastes that have expired
// Returns the number of deleted pastes and any error encountered
func CleanupExpiredPastes() (int, error) {
	// Get current time
	now := time.Now()
	log.Printf("CleanupExpiredPastes started at: %v", now)

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	log.Println("Transaction started successfully.")

	// 1. Lấy ID các paste hết hạn
	rows, err := tx.Query(`SELECT id FROM paste WHERE expires_at IS NOT NULL AND expires_at < ?`, now)
	if err != nil {
		log.Printf("Error fetching expired pastes: %v", err)
		return 0, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning paste ID: %v", err)
			return 0, err
		}
		ids = append(ids, id)
	}

	log.Printf("Found %d expired pastes to delete.", len(ids))

	// 2. Xóa paste_views
	for _, id := range ids {
		if _, err := tx.Exec(`DELETE FROM paste_views WHERE paste_id = ?`, id); err != nil {
			log.Printf("Error deleting paste_views for paste_id: %v, error: %v", id, err)
			return 0, err
		}
		log.Printf("Deleted paste_views for paste_id: %v", id)
	}

	// 3. Xóa paste
	res, err := tx.Exec(`DELETE FROM paste WHERE expires_at IS NOT NULL AND expires_at < ?`, now)
	if err != nil {
		log.Printf("Error deleting expired pastes from paste table: %v", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting number of affected rows: %v", err)
		return 0, err
	}

	log.Printf("Number of expired pastes deleted: %d", rowsAffected)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return 0, err
	}

	log.Println("Transaction committed successfully.")
	return int(rowsAffected), nil
}
