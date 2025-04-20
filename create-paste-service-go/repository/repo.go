package repository

import (
	"database/sql"
	"errors"
	"log"

	"create-paste-service-go/db"
	"create-paste-service-go/model"
)

// ExistsById checks if a paste with the given ID exists
func ExistsById(id string) (bool, error) {
	query := "SELECT id FROM paste WHERE id = ?"

	var pasteId string
	row := db.DB.QueryRow(query, id).Scan(&pasteId)

	if errors.Is(row, sql.ErrNoRows) {
		return false, nil
	}

	log.Printf("Error checking if paste exists: %v", pasteId)

	return true, nil
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
