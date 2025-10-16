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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testCustomPageID       = "550e8400-e29b-41d4-a716-446655440000"
	testCustomPageAuthorID = "550e8400-e29b-41d4-a716-446655440001"
	testCustomPageURL      = "/about-us"
)

// MockCustomPageUseCase is a mock implementation of usecase.CustomPage.
type MockCustomPageUseCase struct {
	mock.Mock
}

func (m *MockCustomPageUseCase) Create(ctx context.Context, authorID string, req *dto.CreateCustomPageRequestDTO) (*dto.CustomPageResponseDTO, error) {
	args := m.Called(ctx, authorID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.CustomPageResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCustomPageUseCase) GetByID(ctx context.Context, id string) (*dto.CustomPageResponseDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.CustomPageResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCustomPageUseCase) GetAll(ctx context.Context) ([]dto.CustomPageResponseDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).([]dto.CustomPageResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockCustomPageUseCase) Update(ctx context.Context, id string, req *dto.UpdateCustomPageRequestDTO) error {
	args := m.Called(ctx, id, req)

	return args.Error(0)
}

func (m *MockCustomPageUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestCustomPageRoutes_GetAll(t *testing.T) {
	t.Run("success - get all custom pages", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.GET("/pages", customPageRouter.GetAll)

		now := time.Now()
		expectedPages := []dto.CustomPageResponseDTO{
			{
				ID:        testCustomPageID,
				CustomURL: testCustomPageURL,
				Content:   "About us content",
				AuthorID:  testCustomPageAuthorID,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Mock expectations
		mockCustomPageUseCase.On("GetAll", mock.Anything).Return(expectedPages, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/pages", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.GET("/pages", customPageRouter.GetAll)

		// Mock expectations
		mockCustomPageUseCase.On("GetAll", mock.Anything).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodGet, "/pages", http.NoBody)
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

		mockCustomPageUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCustomPageRoutes_GetByID(t *testing.T) {
	t.Run("success - get custom page by id", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.GET("/pages/:id", customPageRouter.GetByID)

		now := time.Now()
		expectedPage := &dto.CustomPageResponseDTO{
			ID:        testCustomPageID,
			CustomURL: testCustomPageURL,
			Content:   "About us content",
			AuthorID:  testCustomPageAuthorID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Mock expectations
		mockCustomPageUseCase.On("GetByID", mock.Anything, testCustomPageID).Return(expectedPage, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/pages/"+testCustomPageID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - page not found", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.GET("/pages/:id", customPageRouter.GetByID)

		// Mock expectations
		mockCustomPageUseCase.On("GetByID", mock.Anything, "non-existent-id").Return(nil, apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/pages/non-existent-id", http.NoBody)
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
		assert.Equal(t, "Page not found", meta["message"])

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.GET("/pages/:id", customPageRouter.GetByID)

		// Mock expectations
		mockCustomPageUseCase.On("GetByID", mock.Anything, testCustomPageID).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodGet, "/pages/"+testCustomPageID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCustomPageUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCustomPageRoutes_Create(t *testing.T) {
	t.Run("success - create custom page", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		// Simulate authenticated user middleware setting user_id
		router.POST("/pages", func(c *gin.Context) {
			c.Set("user_id", testCustomPageAuthorID)
			customPageRouter.Create(c)
		})

		requestBody := map[string]string{
			"custom_url": testCustomPageURL,
			"content":    "About us content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		now := time.Now()
		expectedPage := &dto.CustomPageResponseDTO{
			ID:        testCustomPageID,
			CustomURL: testCustomPageURL,
			Content:   "About us content",
			AuthorID:  testCustomPageAuthorID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Mock expectations
		mockCustomPageUseCase.On("Create", mock.Anything, testCustomPageAuthorID, &dto.CreateCustomPageRequestDTO{
			CustomURL: testCustomPageURL,
			Content:   "About us content",
		}).Return(expectedPage, nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(bodyBytes))
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

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.POST("/pages", func(c *gin.Context) {
			c.Set("user_id", testCustomPageAuthorID)
			customPageRouter.Create(c)
		})

		// Missing required fields
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(bodyBytes))
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

		mockCustomPageUseCase.AssertNotCalled(t, "Create")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - user not authenticated", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.POST("/pages", customPageRouter.Create)

		requestBody := map[string]string{
			"custom_url": testCustomPageURL,
			"content":    "About us content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(401), meta["code"])
		assert.Equal(t, "User not authenticated", meta["message"])

		mockCustomPageUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("error - invalid user id format", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.POST("/pages", func(c *gin.Context) {
			c.Set("user_id", 12345) // Set as int instead of string
			customPageRouter.Create(c)
		})

		requestBody := map[string]string{
			"custom_url": testCustomPageURL,
			"content":    "About us content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(500), meta["code"])
		assert.Equal(t, "Invalid user ID format", meta["message"])

		mockCustomPageUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.POST("/pages", func(c *gin.Context) {
			c.Set("user_id", testCustomPageAuthorID)
			customPageRouter.Create(c)
		})

		requestBody := map[string]string{
			"custom_url": testCustomPageURL,
			"content":    "About us content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCustomPageUseCase.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCustomPageUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCustomPageRoutes_Update(t *testing.T) {
	t.Run("success - update custom page", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.PUT("/pages/:id", customPageRouter.Update)

		requestBody := map[string]string{
			"custom_url": "/about-company",
			"content":    "Updated content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCustomPageUseCase.On("Update", mock.Anything, testCustomPageID, &dto.UpdateCustomPageRequestDTO{
			CustomURL: "/about-company",
			Content:   "Updated content",
		}).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/pages/"+testCustomPageID, bytes.NewBuffer(bodyBytes))
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

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.PUT("/pages/:id", customPageRouter.Update)

		// Missing required fields
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/pages/"+testCustomPageID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockCustomPageUseCase.AssertNotCalled(t, "Update")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - page not found", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.PUT("/pages/:id", customPageRouter.Update)

		requestBody := map[string]string{
			"custom_url": "/about-company",
			"content":    "Updated content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCustomPageUseCase.On("Update", mock.Anything, "non-existent-id", mock.Anything).Return(apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/pages/non-existent-id", bytes.NewBuffer(bodyBytes))
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
		assert.Equal(t, "Page not found", meta["message"])

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.PUT("/pages/:id", customPageRouter.Update)

		requestBody := map[string]string{
			"custom_url": "/about-company",
			"content":    "Updated content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCustomPageUseCase.On("Update", mock.Anything, testCustomPageID, mock.Anything).Return(apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/pages/"+testCustomPageID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCustomPageUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestCustomPageRoutes_Delete(t *testing.T) {
	t.Run("success - delete custom page", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.DELETE("/pages/:id", customPageRouter.Delete)

		// Mock expectations
		mockCustomPageUseCase.On("Delete", mock.Anything, testCustomPageID).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/pages/"+testCustomPageID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - page not found", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.DELETE("/pages/:id", customPageRouter.Delete)

		// Mock expectations
		mockCustomPageUseCase.On("Delete", mock.Anything, "non-existent-id").Return(apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/pages/non-existent-id", http.NoBody)
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
		assert.Equal(t, "Page not found", meta["message"])

		mockCustomPageUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCustomPageUseCase := new(MockCustomPageUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		customPageRouter := &customPageRoutes{
			customPage: mockCustomPageUseCase,
			log:        mockLogger,
		}

		router.DELETE("/pages/:id", customPageRouter.Delete)

		// Mock expectations
		mockCustomPageUseCase.On("Delete", mock.Anything, testCustomPageID).Return(apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/pages/"+testCustomPageID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockCustomPageUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
