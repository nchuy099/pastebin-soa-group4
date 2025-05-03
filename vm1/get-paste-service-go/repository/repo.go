package repository

import (
	"database/sql"
	"errors"
	"log"

	"get-paste-service/db"
	"get-paste-service/model"
)

var ErrPasteExpired = errors.New("paste has expired")

func GetPasteByID(id string) (*model.Paste, error) {
	paste := &model.Paste{}

	query := `
    SELECT id, content, title, language, created_at, expires_at, visibility
    FROM paste WHERE id = ? AND (expires_at IS NULL OR expires_at > NOW())
`
	err := db.DB.QueryRow(query, id).Scan(&paste.ID, &paste.Content, &paste.Title, &paste.Language,
		&paste.CreatedAt, &paste.ExpiresAt, &paste.Visibility)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Printf("Error getting paste from database: %v", err.Error())
		return nil, err
	}

	return paste, nil
}
