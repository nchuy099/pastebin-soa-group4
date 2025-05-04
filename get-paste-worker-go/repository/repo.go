package repository

import (
	"log"

	"get-paste-worker/db"
	"get-paste-worker/model"
)

// AddPasteViewDirect adds a paste view record directly to the database
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
