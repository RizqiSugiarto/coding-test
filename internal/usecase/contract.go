package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
)

type Auth interface {
	Login(ctx context.Context, req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
}

type Category interface {
	Create(ctx context.Context, req *dto.CreateCategoryRequestDTO) (*dto.CategoryResponseDTO, error)
	GetByID(ctx context.Context, id string) (*dto.CategoryResponseDTO, error)
	GetAll(ctx context.Context) ([]dto.CategoryResponseDTO, error)
	Update(ctx context.Context, id string, req *dto.UpdateCategoryRequestDTO) error
	Delete(ctx context.Context, id string) error
}
