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

//nolint:dupl // The News and CustomPage interfaces are conceptually different, duplication is intentional
type News interface {
	Create(ctx context.Context, authorID string, req *dto.CreateNewsRequestDTO) (*dto.NewsResponseDTO, error)
	GetByID(ctx context.Context, id string) (*dto.NewsResponseDTO, error)
	GetAll(ctx context.Context) ([]dto.NewsResponseDTO, error)
	Update(ctx context.Context, id string, req *dto.UpdateNewsRequestDTO) error
	Delete(ctx context.Context, id string) error
}

//nolint:dupl // The News and CustomPage interfaces are conceptually different, duplication is intentional
type CustomPage interface {
	Create(ctx context.Context, authorID string, req *dto.CreateCustomPageRequestDTO) (*dto.CustomPageResponseDTO, error)
	GetByID(ctx context.Context, id string) (*dto.CustomPageResponseDTO, error)
	GetAll(ctx context.Context) ([]dto.CustomPageResponseDTO, error)
	Update(ctx context.Context, id string, req *dto.UpdateCustomPageRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type Comment interface {
	Create(ctx context.Context, req *dto.CreateCommentRequestDTO) error
}
