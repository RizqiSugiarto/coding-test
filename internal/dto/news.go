package dto

import "time"

// CreateNewsRequestDTO represents the request to create news.
type CreateNewsRequestDTO struct {
	CategoryID string `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// UpdateNewsRequestDTO represents the request to update news.
type UpdateNewsRequestDTO struct {
	CategoryID string `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// NewsResponseDTO represents the news response.
type NewsResponseDTO struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	AuthorID   string    `json:"author_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
