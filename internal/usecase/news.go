package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/internal/repository"
)

type NewsUseCase struct {
	newsRepo repository.NewsRepo
}

func NewNewsUseCase(newsRepo repository.NewsRepo) *NewsUseCase {
	return &NewsUseCase{
		newsRepo: newsRepo,
	}
}

func (nu *NewsUseCase) Create(ctx context.Context, authorID string, req *dto.CreateNewsRequestDTO) (*dto.NewsResponseDTO, error) {
	news := &entity.News{
		CategoryID: req.CategoryID,
		AuthorID:   authorID,
		Title:      req.Title,
		Content:    req.Content,
	}

	result, err := nu.newsRepo.Create(ctx, news)
	if err != nil {
		return nil, err
	}

	return &dto.NewsResponseDTO{
		ID:         result.ID,
		CategoryID: result.CategoryID,
		AuthorID:   result.AuthorID,
		Title:      result.Title,
		Content:    result.Content,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}, nil
}

func (nu *NewsUseCase) GetByID(ctx context.Context, id string) (*dto.NewsResponseDTO, error) {
	news, err := nu.newsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.NewsResponseDTO{
		ID:         news.ID,
		CategoryID: news.CategoryID,
		AuthorID:   news.AuthorID,
		Title:      news.Title,
		Content:    news.Content,
		CreatedAt:  news.CreatedAt,
		UpdatedAt:  news.UpdatedAt,
	}, nil
}

func (nu *NewsUseCase) GetAll(ctx context.Context) ([]dto.NewsResponseDTO, error) {
	newsList, err := nu.newsRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.NewsResponseDTO, 0, len(newsList))

	for i := range newsList {
		result = append(result, dto.NewsResponseDTO{
			ID:         newsList[i].ID,
			CategoryID: newsList[i].CategoryID,
			AuthorID:   newsList[i].AuthorID,
			Title:      newsList[i].Title,
			Content:    newsList[i].Content,
			CreatedAt:  newsList[i].CreatedAt,
			UpdatedAt:  newsList[i].UpdatedAt,
		})
	}

	return result, nil
}

func (nu *NewsUseCase) Update(ctx context.Context, id string, req *dto.UpdateNewsRequestDTO) error {
	news := &entity.News{
		ID:         id,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
	}

	err := nu.newsRepo.Update(ctx, news)
	if err != nil {
		return err
	}

	return nil
}

func (nu *NewsUseCase) Delete(ctx context.Context, id string) error {
	err := nu.newsRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
