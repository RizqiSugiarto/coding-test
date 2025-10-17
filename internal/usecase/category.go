package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/internal/repository"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepo
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (cu *CategoryUseCase) Create(ctx context.Context, req *dto.CreateCategoryRequestDTO) (*dto.CategoryResponseDTO, error) {
	category := &entity.Category{
		Name: req.Name,
	}

	result, err := cu.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponseDTO{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (cu *CategoryUseCase) GetByID(ctx context.Context, id string) (*dto.CategoryResponseDTO, error) {
	category, err := cu.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponseDTO{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (cu *CategoryUseCase) GetAll(ctx context.Context) ([]dto.CategoryResponseDTO, error) {
	categories, err := cu.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CategoryResponseDTO, 0, len(categories))

	for _, category := range categories {
		result = append(result, dto.CategoryResponseDTO{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
	}

	return result, nil
}

func (cu *CategoryUseCase) Update(ctx context.Context, id string, req *dto.UpdateCategoryRequestDTO) error {
	category := &entity.Category{
		ID:   id,
		Name: req.Name,
	}

	err := cu.categoryRepo.Update(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (cu *CategoryUseCase) Delete(ctx context.Context, id string) error {
	err := cu.categoryRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
