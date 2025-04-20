package models

import (
	"database/sql"
	"time"
)

type Paste struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Title      string     `json:"title"`
	Language   string     `json:"language"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Views      int        `json:"views"`
	Visibility string     `json:"visibility"`
	Status     string     `json:"status"`
}

func GetPublicPastes() ([]Paste, error) {
	query := `
		SELECT id, content, title, language, created_at, expires_at, views, visibility, status
		FROM pastes
		WHERE visibility = 'public' AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC
		LIMIT 10`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var p Paste
		var expires sql.NullTime
		if err := rows.Scan(&p.ID, &p.Content, &p.Title, &p.Language, &p.CreatedAt, &expires, &p.Views, &p.Visibility, &p.Status); err != nil {
			return nil, err
		}
		if expires.Valid {
			p.ExpiresAt = &expires.Time
		}
		pastes = append(pastes, p)
	}
	return pastes, nil
}
