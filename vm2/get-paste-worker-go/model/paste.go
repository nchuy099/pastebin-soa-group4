package model

import "time"

// PasteViews represents a view record for a paste
type PasteViews struct {
	PasteID  string    `json:"paste_id"`
	ViewedAt time.Time `json:"viewed_at"`
}
