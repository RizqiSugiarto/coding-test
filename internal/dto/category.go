package dto

import "time"

type CreateCategoryRequestDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequestDTO struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
