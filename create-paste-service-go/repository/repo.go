package repository

import (
	"database/sql"
	"log"

	"create-paste-service-go/db"
	"create-paste-service-go/model"
)

// ExistsById checks if a paste with the given ID exists
func ExistsById(id string) (bool, error) {
	query := "SELECT 1 FROM paste WHERE id = ? LIMIT 1"
	var exists bool
	err := db.DB.QueryRow(query, id).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking if paste exists: %v", err.Error())
		return false, err
	}
	return exists, nil
}

// SavePaste saves a new paste to the database
func SavePaste(paste *model.Paste) error {
	query := `
		INSERT INTO paste (id, content, title, language, created_at, expires_at, views, visibility)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.DB.Exec(
		query,
		paste.ID,
		paste.Content,
		paste.Title,
		paste.Language,
		paste.CreatedAt,
		paste.ExpiresAt,
		paste.Views,
		paste.Visibility,
	)

	if err != nil {
		log.Printf("Error saving paste: %v", err.Error())
		return err
	}

	return nil
}
