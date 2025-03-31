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
