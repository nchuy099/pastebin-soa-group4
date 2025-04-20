package model

import "time"

type Paste struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Title      string     `json:"title"`
	Language   string     `json:"language"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Visibility string     `json:"visibility"`
}

type PasteView struct {
	PasteID  string    `json:"paste_id"`
	ViewedAt time.Time `json:"viewed_at"`
}

// ResponseData represents a standard API response format
type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}
