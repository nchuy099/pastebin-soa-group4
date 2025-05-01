package scheduler

import (
	"cleanup-expired-service-go/db"
	"log"
	"time"
)

// CleanupExpiredPastes deletes all pastes that have expired
// Returns the number of deleted pastes and any error encountered
func CleanupExpiredPastes() (int, error) {
	now := time.Now()

	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return 0, err
	}

	// 1. Delete from paste_view first
	deleteViewsQuery := `
		DELETE FROM paste_views
		WHERE paste_id IN (
			SELECT id FROM paste WHERE expires_at IS NOT NULL AND expires_at < ?
		)`
	if _, err := tx.Exec(deleteViewsQuery, now); err != nil {
		log.Printf("Error deleting from paste_views: %v", err)
		tx.Rollback()
		return 0, err
	}

	// 2. Delete from paste
	deletePastesQuery := `
		DELETE FROM paste
		WHERE expires_at IS NOT NULL AND expires_at < ?`
	result, err := tx.Exec(deletePastesQuery, now)
	if err != nil {
		log.Printf("Error deleting from paste: %v", err)
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return 0, err
	}

	// Get number of deleted pastes
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows: %v", err)
		return 0, err
	}

	return int(rowsAffected), nil
}
