package model

import (
	"time"
)

// Visibility type defines visibility options for pastes
type Visibility string

// Visibility constants
const (
	Public   Visibility = "PUBLIC"
	Unlisted Visibility = "UNLISTED"
)

// Paste represents a paste entity
type Paste struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Title      string     `json:"title"`
	Language   string     `json:"language"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Visibility Visibility `json:"visibility"`
}

// CreatePasteRequest represents a request to create a new paste
type CreatePasteRequest struct {
	Content    string `json:"content"`
	Title      string `json:"title"`
	Language   string `json:"language"`
	ExpiresIn  int64  `json:"expiresIn"`
	Visibility string `json:"visibility"`
}

// ResponseData represents a standard API response format
type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}
