package repository

import (
	"log"

	"get-public-service/db"
	"get-public-service/model"
)

func GetPublicPastes(limit, offset int) ([]model.Paste, error) {
	query := `SELECT id, content, title, language, created_at, expires_at, views, visibility 
	          FROM paste 
	          WHERE visibility = 'public' AND (expires_at IS NULL OR expires_at > NOW()) 
	          ORDER BY created_at DESC 
	          LIMIT ? OFFSET ?`

	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error querying public pastes:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var pastes []model.Paste
	for rows.Next() {
		var paste model.Paste
		err := rows.Scan(&paste.ID, &paste.Content, &paste.Title, &paste.Language,
			&paste.CreatedAt, &paste.ExpiresAt, &paste.Views, &paste.Visibility)
		if err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}

		pastes = append(pastes, paste)
	}

	return pastes, nil
}

// CountPublicPastes returns the total number of public pastes
func CountPublicPastes() (int, error) {
	query := `SELECT COUNT(*) FROM paste 
	          WHERE visibility = 'public' AND (expires_at IS NULL OR expires_at > NOW())`

	var count int
	err := db.DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println("Error counting public pastes:", err.Error())
		return 0, err
	}

	return count, nil
}
