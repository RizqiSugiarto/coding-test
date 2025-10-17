package dto

import "time"

// CreateCustomPageRequestDTO represents the request to create a custom page.
type CreateCustomPageRequestDTO struct {
	CustomURL string `json:"custom_url" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// UpdateCustomPageRequestDTO represents the request to update a custom page.
type UpdateCustomPageRequestDTO struct {
	CustomURL string `json:"custom_url" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// CustomPageResponseDTO represents the response for a custom page.
type CustomPageResponseDTO struct {
	ID        string    `json:"id"`
	CustomURL string    `json:"custom_url"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
