package repository

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type CategoryRepo interface {
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	GetAll(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
}
