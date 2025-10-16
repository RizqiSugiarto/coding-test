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
	testPageID           = "550e8400-e29b-41d4-a716-446655440000"
	testPageAuthorID     = "550e8400-e29b-41d4-a716-446655440001"
	nonExistentPageID    = "non-existent-id"
	testPageCustomURL    = "/about-us"
	testPageCustomURLNew = "/about-company"
	testPageContent      = "Updated content"
)

// MockCustomPageRepo is a mock implementation of repository.CustomPageRepo.
type MockCustomPageRepo struct {
	mock.Mock
}

func (m *MockCustomPageRepo) Create(ctx context.Context, page *entity.CustomPage) (*entity.CustomPage, error) {
	args := m.Called(ctx, page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.CustomPage)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCustomPageRepo) GetByID(ctx context.Context, id string) (*entity.CustomPage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.CustomPage)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCustomPageRepo) GetAll(ctx context.Context) ([]entity.CustomPage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).([]entity.CustomPage)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCustomPageRepo) Update(ctx context.Context, page *entity.CustomPage) error {
	args := m.Called(ctx, page)

	return args.Error(0)
}

func (m *MockCustomPageRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestCustomPageUseCase_Create(t *testing.T) {
	t.Run("success - create custom page", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCustomPageRequestDTO{
			CustomURL: testPageCustomURL,
			Content:   "This is the about us page content",
		}

		now := time.Now()
		expectedPage := &entity.CustomPage{
			ID:        testPageID,
			CustomURL: testPageCustomURL,
			Content:   "This is the about us page content",
			AuthorID:  testPageAuthorID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(page *entity.CustomPage) bool {
			return page.CustomURL == testPageCustomURL &&
				page.AuthorID == testPageAuthorID &&
				page.Content == "This is the about us page content"
		})).Return(expectedPage, nil)

		result, err := useCase.Create(ctx, testPageAuthorID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedPage.ID, result.ID)
		assert.Equal(t, expectedPage.CustomURL, result.CustomURL)
		assert.Equal(t, expectedPage.Content, result.Content)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository create fails", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCustomPageRequestDTO{
			CustomURL: testPageCustomURL,
			Content:   "This is the about us page content",
		}

		mockRepo.On("Create", ctx, mock.Anything).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.Create(ctx, testPageAuthorID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomPageUseCase_GetByID(t *testing.T) {
	t.Run("success - get custom page by id", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		now := time.Now()
		expectedPage := &entity.CustomPage{
			ID:        testPageID,
			CustomURL: testPageCustomURL,
			Content:   "This is the about us page content",
			AuthorID:  testPageAuthorID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.On("GetByID", ctx, testPageID).Return(expectedPage, nil)

		result, err := useCase.GetByID(ctx, testPageID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedPage.ID, result.ID)
		assert.Equal(t, expectedPage.CustomURL, result.CustomURL)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - custom page not found", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		pageID := nonExistentPageID

		mockRepo.On("GetByID", ctx, pageID).Return(nil, apperror.ErrNotFound)

		result, err := useCase.GetByID(ctx, pageID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository get fails", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("GetByID", ctx, testPageID).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.GetByID(ctx, testPageID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomPageUseCase_GetAll(t *testing.T) {
	t.Run("success - get all custom pages", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		now := time.Now()
		pageList := []entity.CustomPage{
			{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				CustomURL: "/about-us",
				Content:   "About content",
				AuthorID:  testPageAuthorID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "550e8400-e29b-41d4-a716-446655440002",
				CustomURL: "/contact",
				Content:   "Contact content",
				AuthorID:  testPageAuthorID,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mockRepo.On("GetAll", ctx).Return(pageList, nil)

		result, err := useCase.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "/about-us", result[0].CustomURL)
		assert.Equal(t, "/contact", result[1].CustomURL)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - get all custom pages empty result", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		pageList := []entity.CustomPage{}

		mockRepo.On("GetAll", ctx).Return(pageList, nil)

		result, err := useCase.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository getall fails", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("GetAll", ctx).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomPageUseCase_Update(t *testing.T) {
	t.Run("success - update custom page", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.UpdateCustomPageRequestDTO{
			CustomURL: testPageCustomURLNew,
			Content:   testPageContent,
		}

		mockRepo.On("Update", ctx, mock.MatchedBy(func(page *entity.CustomPage) bool {
			return page.ID == testPageID &&
				page.CustomURL == testPageCustomURLNew &&
				page.Content == testPageContent
		})).Return(nil)

		err := useCase.Update(ctx, testPageID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - custom page not found", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		pageID := nonExistentPageID
		req := &dto.UpdateCustomPageRequestDTO{
			CustomURL: testPageCustomURLNew,
			Content:   testPageContent,
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(apperror.ErrNotFound)

		err := useCase.Update(ctx, pageID, req)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository update fails", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.UpdateCustomPageRequestDTO{
			CustomURL: testPageCustomURLNew,
			Content:   testPageContent,
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(apperror.ErrDatabaseConnection)

		err := useCase.Update(ctx, testPageID, req)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomPageUseCase_Delete(t *testing.T) {
	t.Run("success - delete custom page", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("Delete", ctx, testPageID).Return(nil)

		err := useCase.Delete(ctx, testPageID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - custom page not found", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()
		pageID := nonExistentPageID

		mockRepo.On("Delete", ctx, pageID).Return(apperror.ErrNotFound)

		err := useCase.Delete(ctx, pageID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository delete fails", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)
		useCase := NewCustomPageUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("Delete", ctx, testPageID).Return(apperror.ErrDatabaseConnection)

		err := useCase.Delete(ctx, testPageID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewCustomPageUseCase(t *testing.T) {
	t.Run("success - create new custom page usecase", func(t *testing.T) {
		mockRepo := new(MockCustomPageRepo)

		useCase := NewCustomPageUseCase(mockRepo)

		assert.NotNil(t, useCase)
		assert.NotNil(t, useCase.customPageRepo)
	})
}
