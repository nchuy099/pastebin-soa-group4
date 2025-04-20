package repository

import (
	"database/sql"
	"errors"
	"log"

	"get-paste-service/cache"
	"get-paste-service/db"
	"get-paste-service/model"
	"get-paste-service/queue"
)

var ErrPasteExpired = errors.New("paste has expired")

func GetPasteByID(id string) (*model.Paste, error) {
	// Try to get paste from cache first
	paste, err := cache.GetPasteFromCache(id)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	// If found in cache, return it
	if paste != nil {
		log.Printf("Cache hit for paste ID: %s", id)
		return paste, nil
	}

	// If not in cache, get from database
	log.Printf("Cache miss for paste ID: %s, fetching from database", id)

	paste = &model.Paste{}
	var expiresAt sql.NullTime

	query := `
    SELECT id, content, title, language, created_at, expires_at, visibility
    FROM paste WHERE id = ? AND (expires_at IS NULL OR expires_at > NOW())
`
	err = db.DB.QueryRow(query, id).Scan(&paste.ID, &paste.Content, &paste.Title, &paste.Language,
		&paste.CreatedAt, &expiresAt, &paste.Visibility)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Printf("Error getting paste: %v", err.Error())
		return nil, err
	}

	if expiresAt.Valid {
		paste.ExpiresAt = &expiresAt.Time
	}

	// Store paste in cache for future requests
	if cacheErr := cache.SetPasteToCache(paste); cacheErr != nil {
		log.Printf("Failed to cache paste: %v", cacheErr)
	}

	return paste, nil
}

func AddPasteView(pasteView *model.PasteViews, remoteAddr, userAgent string) error {
	// Publish the view update to RabbitMQ
	_ = queue.PublishViewUpdate(pasteView.PasteID, pasteView.ViewedAt, remoteAddr, userAgent)

	// Asynchronously update database (handled by worker)
	// The database update will be handled by the view-update-worker
	// This allows for faster response times
	return nil
}

// AddPasteViewDirect is used by the worker to directly update the database
func AddPasteViewDirect(pasteView *model.PasteViews) error {
	query := `
		INSERT INTO paste_views (paste_id, viewed_at)
		VALUES (?, ?)
	`
	_, err := db.DB.Exec(query, pasteView.PasteID, pasteView.ViewedAt)
	if err != nil {
		log.Printf("Error adding paste view: %v", err.Error())
	}
	return nil
}
