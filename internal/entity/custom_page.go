package entity

import "time"

// CustomPage represents a custom page in the system.
type CustomPage struct {
	ID        string    `json:"id"`
	CustomURL string    `json:"custom_url"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
