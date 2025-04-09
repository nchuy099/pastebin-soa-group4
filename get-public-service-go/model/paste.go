package model

import "time"

type Paste struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Title      string     `json:"title"`
	Language   string     `json:"language"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Views      int        `json:"views"`
	Visibility string     `json:"visibility"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems,omitempty"`
}

// PasteListResponse represents a paginated list of pastes
type PasteListResponse struct {
	Pastes     []Paste    `json:"pastes"`
	Pagination Pagination `json:"pagination"`
}

// ResponseData represents a standard API response format
type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
