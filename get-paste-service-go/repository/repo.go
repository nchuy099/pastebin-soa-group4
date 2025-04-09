package repository

import (
	"database/sql"
	"errors"
	"time"

	"get-paste-service/model"
	"get-paste-service/db"
)

var ErrPasteExpired = errors.New("paste has expired")

func GetPasteByID(id string) (*model.Paste, error) {
	var paste model.Paste
	var expiresAt sql.NullTime

	query := `SELECT id, content, title, language, created_at, expires_at, views, visibility 
	          FROM paste WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(&paste.ID, &paste.Content, &paste.Title, &paste.Language,
		&paste.CreatedAt, &expiresAt, &paste.Views, &paste.Visibility)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	if expiresAt.Valid {
		paste.ExpiresAt = &expiresAt.Time
		if expiresAt.Time.Before(time.Now()) {
			return nil, ErrPasteExpired
		}
	}

	// Cập nhật lượt xem (Không blocking API)
	go func() {
		_, _ = db.DB.Exec("UPDATE paste SET views = views + 1 WHERE id = ?", id)
	}()

	return &paste, nil
}
