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
	testCategoryID = "550e8400-e29b-41d4-a716-446655440000"
	nonExistentID  = "non-existent-id"
)

// MockCategoryRepo is a mock implementation of repository.CategoryRepo.
type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.Category)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCategoryRepo) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.Category)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCategoryRepo) GetAll(ctx context.Context) ([]entity.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).([]entity.Category)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCategoryRepo) Update(ctx context.Context, category *entity.Category) error {
	args := m.Called(ctx, category)

	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestCategoryUseCase_Create(t *testing.T) {
	t.Run("success - create category", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCategoryRequestDTO{
			Name: "Technology",
		}

		now := time.Now()
		expectedCategory := &entity.Category{
			ID:        testCategoryID,
			Name:      "Technology",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(cat *entity.Category) bool {
			return cat.Name == "Technology"
		})).Return(expectedCategory, nil)

		result, err := useCase.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCategory.ID, result.ID)
		assert.Equal(t, expectedCategory.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository create fails", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCategoryRequestDTO{
			Name: "Technology",
		}

		mockRepo.On("Create", ctx, mock.Anything).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUseCase_GetByID(t *testing.T) {
	t.Run("success - get category by id", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := testCategoryID

		now := time.Now()
		expectedCategory := &entity.Category{
			ID:        categoryID,
			Name:      "Technology",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

		result, err := useCase.GetByID(ctx, categoryID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCategory.ID, result.ID)
		assert.Equal(t, expectedCategory.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := nonExistentID

		mockRepo.On("GetByID", ctx, categoryID).Return(nil, apperror.ErrNotFound)

		result, err := useCase.GetByID(ctx, categoryID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository get fails", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := testCategoryID

		mockRepo.On("GetByID", ctx, categoryID).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.GetByID(ctx, categoryID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUseCase_GetAll(t *testing.T) {
	t.Run("success - get all categories", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()

		now := time.Now()
		categories := []entity.Category{
			{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				Name:      "Technology",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "550e8400-e29b-41d4-a716-446655440002",
				Name:      "Sports",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mockRepo.On("GetAll", ctx).Return(categories, nil)

		result, err := useCase.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "Technology", result[0].Name)
		assert.Equal(t, "Sports", result[1].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - get all categories empty result", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()

		categories := []entity.Category{}

		mockRepo.On("GetAll", ctx).Return(categories, nil)

		result, err := useCase.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository getall fails", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()

		mockRepo.On("GetAll", ctx).Return(nil, apperror.ErrDatabaseConnection)

		result, err := useCase.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUseCase_Update(t *testing.T) {
	t.Run("success - update category", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := testCategoryID
		req := &dto.UpdateCategoryRequestDTO{
			Name: "Updated Technology",
		}

		mockRepo.On("Update", ctx, mock.MatchedBy(func(cat *entity.Category) bool {
			return cat.ID == categoryID && cat.Name == "Updated Technology"
		})).Return(nil)

		err := useCase.Update(ctx, categoryID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := nonExistentID
		req := &dto.UpdateCategoryRequestDTO{
			Name: "Updated Technology",
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(apperror.ErrNotFound)

		err := useCase.Update(ctx, categoryID, req)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository update fails", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := testCategoryID
		req := &dto.UpdateCategoryRequestDTO{
			Name: "Updated Technology",
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(apperror.ErrDatabaseConnection)

		err := useCase.Update(ctx, categoryID, req)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUseCase_Delete(t *testing.T) {
	t.Run("success - delete category", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := testCategoryID

		mockRepo.On("Delete", ctx, categoryID).Return(nil)

		err := useCase.Delete(ctx, categoryID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := nonExistentID

		mockRepo.On("Delete", ctx, categoryID).Return(apperror.ErrNotFound)

		err := useCase.Delete(ctx, categoryID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository delete fails", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		useCase := NewCategoryUseCase(mockRepo)

		ctx := context.Background()
		categoryID := testCategoryID

		mockRepo.On("Delete", ctx, categoryID).Return(apperror.ErrDatabaseConnection)

		err := useCase.Delete(ctx, categoryID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewCategoryUseCase(t *testing.T) {
	t.Run("success - create new category usecase", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)

		useCase := NewCategoryUseCase(mockRepo)

		assert.NotNil(t, useCase)
		assert.NotNil(t, useCase.categoryRepo)
	})
}
