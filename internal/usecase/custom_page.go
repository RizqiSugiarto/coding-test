package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/internal/repository"
)

type CustomPageUseCase struct {
	customPageRepo repository.CustomPageRepo
}

func NewCustomPageUseCase(customPageRepo repository.CustomPageRepo) *CustomPageUseCase {
	return &CustomPageUseCase{
		customPageRepo: customPageRepo,
	}
}

func (cu *CustomPageUseCase) Create(ctx context.Context, authorID string, req *dto.CreateCustomPageRequestDTO) (*dto.CustomPageResponseDTO, error) {
	page := &entity.CustomPage{
		CustomURL: req.CustomURL,
		Content:   req.Content,
		AuthorID:  authorID,
	}

	result, err := cu.customPageRepo.Create(ctx, page)
	if err != nil {
		return nil, err
	}

	return &dto.CustomPageResponseDTO{
		ID:        result.ID,
		CustomURL: result.CustomURL,
		Content:   result.Content,
		AuthorID:  result.AuthorID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (cu *CustomPageUseCase) GetByID(ctx context.Context, id string) (*dto.CustomPageResponseDTO, error) {
	page, err := cu.customPageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.CustomPageResponseDTO{
		ID:        page.ID,
		CustomURL: page.CustomURL,
		Content:   page.Content,
		AuthorID:  page.AuthorID,
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
	}, nil
}

func (cu *CustomPageUseCase) GetAll(ctx context.Context) ([]dto.CustomPageResponseDTO, error) {
	pageList, err := cu.customPageRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CustomPageResponseDTO, 0, len(pageList))

	for i := range pageList {
		result = append(result, dto.CustomPageResponseDTO{
			ID:        pageList[i].ID,
			CustomURL: pageList[i].CustomURL,
			Content:   pageList[i].Content,
			AuthorID:  pageList[i].AuthorID,
			CreatedAt: pageList[i].CreatedAt,
			UpdatedAt: pageList[i].UpdatedAt,
		})
	}

	return result, nil
}

func (cu *CustomPageUseCase) Update(ctx context.Context, id string, req *dto.UpdateCustomPageRequestDTO) error {
	page := &entity.CustomPage{
		ID:        id,
		CustomURL: req.CustomURL,
		Content:   req.Content,
	}

	err := cu.customPageRepo.Update(ctx, page)
	if err != nil {
		return err
	}

	return nil
}

func (cu *CustomPageUseCase) Delete(ctx context.Context, id string) error {
	err := cu.customPageRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
