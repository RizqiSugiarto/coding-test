package entity

import "time"

// News represents a news article in the system.
type News struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	AuthorID   string    `json:"author_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
