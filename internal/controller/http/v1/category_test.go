package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testCategoryID = "550e8400-e29b-41d4-a716-446655440000"

// MockCategoryUseCase is a mock implementation of usecase.Category.
type MockCategoryUseCase struct {
	mock.Mock
}

func (m *MockCategoryUseCase) Create(ctx context.Context, req *dto.CreateCategoryRequestDTO) (*dto.CategoryResponseDTO, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.CategoryResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCategoryUseCase) GetByID(ctx context.Context, id string) (*dto.CategoryResponseDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.CategoryResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCategoryUseCase) GetAll(ctx context.Context) ([]dto.CategoryResponseDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).([]dto.CategoryResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCategoryUseCase) Update(ctx context.Context, id string, req *dto.UpdateCategoryRequestDTO) error {
	args := m.Called(ctx, id, req)

	return args.Error(0)
}

func (m *MockCategoryUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestCategoryRoutes_GetAll(t *testing.T) {
	t.Run("success - get all categories", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.GET("/categories", categoryRouter.GetAll)

		now := time.Now()
		expectedCategories := []dto.CategoryResponseDTO{
			{
				ID:        testCategoryID,
				Name:      "Technology",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Mock expectations
		mockCategoryUseCase.On("GetAll", mock.Anything).Return(expectedCategories, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/categories", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.GET("/categories", categoryRouter.GetAll)

		// Mock expectations
		mockCategoryUseCase.On("GetAll", mock.Anything).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodGet, "/categories", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(500), meta["code"])
		assert.Equal(t, "Internal server error", meta["message"])

		mockCategoryUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCategoryRoutes_GetByID(t *testing.T) {
	t.Run("success - get category by id", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.GET("/categories/:id", categoryRouter.GetByID)

		now := time.Now()
		expectedCategory := &dto.CategoryResponseDTO{
			ID:        testCategoryID,
			Name:      "Technology",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Mock expectations
		mockCategoryUseCase.On("GetByID", mock.Anything, testCategoryID).Return(expectedCategory, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/categories/"+testCategoryID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.GET("/categories/:id", categoryRouter.GetByID)

		// Mock expectations
		mockCategoryUseCase.On("GetByID", mock.Anything, "non-existent-id").Return(nil, apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/categories/non-existent-id", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(404), meta["code"])
		assert.Equal(t, "Category not found", meta["message"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.GET("/categories/:id", categoryRouter.GetByID)

		// Mock expectations
		mockCategoryUseCase.On("GetByID", mock.Anything, testCategoryID).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodGet, "/categories/"+testCategoryID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCategoryUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCategoryRoutes_Create(t *testing.T) {
	t.Run("success - create category", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.POST("/categories", categoryRouter.Create)

		requestBody := map[string]string{
			"name": "Technology",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		now := time.Now()
		expectedCategory := &dto.CategoryResponseDTO{
			ID:        testCategoryID,
			Name:      "Technology",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Mock expectations
		mockCategoryUseCase.On("Create", mock.Anything, &dto.CreateCategoryRequestDTO{
			Name: "Technology",
		}).Return(expectedCategory, nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(201), meta["code"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.POST("/categories", categoryRouter.Create)

		// Missing required field
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(400), meta["code"])
		assert.Equal(t, "Invalid request payload", meta["message"])

		mockCategoryUseCase.AssertNotCalled(t, "Create")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.POST("/categories", categoryRouter.Create)

		requestBody := map[string]string{
			"name": "Technology",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCategoryUseCase.On("Create", mock.Anything, mock.Anything).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCategoryUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCategoryRoutes_Update(t *testing.T) {
	t.Run("success - update category", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.PUT("/categories/:id", categoryRouter.Update)

		requestBody := map[string]string{
			"name": "Updated Technology",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCategoryUseCase.On("Update", mock.Anything, testCategoryID, &dto.UpdateCategoryRequestDTO{
			Name: "Updated Technology",
		}).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/categories/"+testCategoryID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.PUT("/categories/:id", categoryRouter.Update)

		// Missing required field
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/categories/"+testCategoryID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockCategoryUseCase.AssertNotCalled(t, "Update")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.PUT("/categories/:id", categoryRouter.Update)

		requestBody := map[string]string{
			"name": "Updated Technology",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCategoryUseCase.On("Update", mock.Anything, "non-existent-id", mock.Anything).Return(apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/categories/non-existent-id", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(404), meta["code"])
		assert.Equal(t, "Category not found", meta["message"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.PUT("/categories/:id", categoryRouter.Update)

		requestBody := map[string]string{
			"name": "Updated Technology",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCategoryUseCase.On("Update", mock.Anything, testCategoryID, mock.Anything).Return(apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/categories/"+testCategoryID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCategoryUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCategoryRoutes_Delete(t *testing.T) {
	t.Run("success - delete category", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.DELETE("/categories/:id", categoryRouter.Delete)

		// Mock expectations
		mockCategoryUseCase.On("Delete", mock.Anything, testCategoryID).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/categories/"+testCategoryID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - category not found", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.DELETE("/categories/:id", categoryRouter.Delete)

		// Mock expectations
		mockCategoryUseCase.On("Delete", mock.Anything, "non-existent-id").Return(apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/categories/non-existent-id", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(404), meta["code"])
		assert.Equal(t, "Category not found", meta["message"])

		mockCategoryUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCategoryUseCase := new(MockCategoryUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		categoryRouter := &categoryRoutes{
			category: mockCategoryUseCase,
			log:      mockLogger,
		}

		router.DELETE("/categories/:id", categoryRouter.Delete)

		// Mock expectations
		mockCategoryUseCase.On("Delete", mock.Anything, testCategoryID).Return(apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/categories/"+testCategoryID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCategoryUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
