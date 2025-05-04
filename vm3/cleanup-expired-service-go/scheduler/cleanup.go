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

	// Query to delete expired pastes
	query := `DELETE FROM paste WHERE expires_at IS NOT NULL AND expires_at < ?`

	// Execute the query
	result, err := db.DB.Exec(query, now)
	if err != nil {
		log.Printf("Error deleting expired pastes: %v", err)
		return 0, err
	}

	// Get the number of affected rows (deleted pastes)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting number of affected rows: %v", err)
		return 0, err
	}

	return int(rowsAffected), nil
}
