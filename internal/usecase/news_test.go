package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testNewsID         = "550e8400-e29b-41d4-a716-446655440000"
	testNewsCategoryID = "550e8400-e29b-41d4-a716-446655440001"
	testNewsAuthorID   = "550e8400-e29b-41d4-a716-446655440002"
	nonExistentNewsID  = "non-existent-id"
)

// MockNewsRepo is a mock implementation of repository.NewsRepo.
type MockNewsRepo struct {
	mock.Mock
}

func (m *MockNewsRepo) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	args := m.Called(ctx, news)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.News)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockNewsRepo) GetByID(ctx context.Context, id string) (*entity.News, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.News)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockNewsRepo) GetAll(ctx context.Context) ([]entity.News, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).([]entity.News)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockNewsRepo) Update(ctx context.Context, news *entity.News) error {
	args := m.Called(ctx, news)

	return args.Error(0)
}

func (m *MockNewsRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestNewsUseCase_Create(t *testing.T) {
	t.Run("success - create news", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Breaking News",
			Content:    "This is the news content",
		}

		now := time.Now()
		expectedNews := &entity.News{
			ID:         testNewsID,
			CategoryID: testNewsCategoryID,
			AuthorID:   testNewsAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(news *entity.News) bool {
			return news.CategoryID == testNewsCategoryID &&
				news.AuthorID == testNewsAuthorID &&
				news.Title == "Breaking News" &&
				news.Content == "This is the news content"
		})).Return(expectedNews, nil)

		result, err := useCase.Create(ctx, testNewsAuthorID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedNews.ID, result.ID)
		assert.Equal(t, expectedNews.Title, result.Title)
		assert.Equal(t, expectedNews.Content, result.Content)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository create fails", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Breaking News",
			Content:    "This is the news content",
		}

		mockRepo.On("Create", ctx, mock.Anything).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.Create(ctx, testNewsAuthorID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewsUseCase_GetByID(t *testing.T) {
	t.Run("success - get news by id", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		now := time.Now()
		expectedNews := &entity.News{
			ID:         testNewsID,
			CategoryID: testNewsCategoryID,
			AuthorID:   testNewsAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.On("GetByID", ctx, testNewsID).Return(expectedNews, nil)

		result, err := useCase.GetByID(ctx, testNewsID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedNews.ID, result.ID)
		assert.Equal(t, expectedNews.Title, result.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - news not found", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		newsID := nonExistentNewsID

		mockRepo.On("GetByID", ctx, newsID).Return(nil, apperror.ErrNotFound)

		result, err := useCase.GetByID(ctx, newsID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository get fails", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("GetByID", ctx, testNewsID).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.GetByID(ctx, testNewsID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewsUseCase_GetAll(t *testing.T) {
	t.Run("success - get all news", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		now := time.Now()
		newsList := []entity.News{
			{
				ID:         "550e8400-e29b-41d4-a716-446655440001",
				CategoryID: testNewsCategoryID,
				AuthorID:   testNewsAuthorID,
				Title:      "News 1",
				Content:    "Content 1",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				ID:         "550e8400-e29b-41d4-a716-446655440002",
				CategoryID: testNewsCategoryID,
				AuthorID:   testNewsAuthorID,
				Title:      "News 2",
				Content:    "Content 2",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}

		mockRepo.On("GetAll", ctx).Return(newsList, nil)

		result, err := useCase.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "News 1", result[0].Title)
		assert.Equal(t, "News 2", result[1].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - get all news empty result", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		newsList := []entity.News{}

		mockRepo.On("GetAll", ctx).Return(newsList, nil)

		result, err := useCase.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository getall fails", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("GetAll", ctx).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewsUseCase_Update(t *testing.T) {
	t.Run("success - update news", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.UpdateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}

		mockRepo.On("Update", ctx, mock.MatchedBy(func(news *entity.News) bool {
			return news.ID == testNewsID &&
				news.CategoryID == testNewsCategoryID &&
				news.Title == "Updated News" &&
				news.Content == "Updated content"
		})).Return(nil)

		err := useCase.Update(ctx, testNewsID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - news not found", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		newsID := nonExistentNewsID
		req := &dto.UpdateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(apperror.ErrNotFound)

		err := useCase.Update(ctx, newsID, req)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository update fails", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.UpdateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(apperror.ErrDatabaseConnection)

		err := useCase.Update(ctx, testNewsID, req)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewsUseCase_Delete(t *testing.T) {
	t.Run("success - delete news", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("Delete", ctx, testNewsID).Return(nil)

		err := useCase.Delete(ctx, testNewsID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - news not found", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()
		newsID := nonExistentNewsID

		mockRepo.On("Delete", ctx, newsID).Return(apperror.ErrNotFound)

		err := useCase.Delete(ctx, newsID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository delete fails", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)
		useCase := NewNewsUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("Delete", ctx, testNewsID).Return(apperror.ErrDatabaseConnection)

		err := useCase.Delete(ctx, testNewsID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewNewsUseCase(t *testing.T) {
	t.Run("success - create new news usecase", func(t *testing.T) {
		mockRepo := new(MockNewsRepo)

		useCase := NewNewsUseCase(mockRepo)

		assert.NotNil(t, useCase)
		assert.NotNil(t, useCase.newsRepo)
	})
}
