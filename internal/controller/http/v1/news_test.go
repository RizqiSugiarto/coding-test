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
	testNewsID         = "550e8400-e29b-41d4-a716-446655440000"
	testNewsAuthorID   = "550e8400-e29b-41d4-a716-446655440001"
	testNewsCategoryID = "550e8400-e29b-41d4-a716-446655440002"
)

// MockNewsUseCase is a mock implementation of usecase.News.
type MockNewsUseCase struct {
	mock.Mock
}

func (m *MockNewsUseCase) Create(ctx context.Context, authorID string, req *dto.CreateNewsRequestDTO) (*dto.NewsResponseDTO, error) {
	args := m.Called(ctx, authorID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.NewsResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockNewsUseCase) GetByID(ctx context.Context, id string) (*dto.NewsResponseDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.NewsResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockNewsUseCase) GetAll(ctx context.Context) ([]dto.NewsResponseDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).([]dto.NewsResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func (m *MockNewsUseCase) Update(ctx context.Context, id string, req *dto.UpdateNewsRequestDTO) error {
	args := m.Called(ctx, id, req)

	return args.Error(0)
}

func (m *MockNewsUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestNewsRoutes_GetAll(t *testing.T) {
	t.Run("success - get all news", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.GET("/news", newsRouter.GetAll)

		now := time.Now()
		expectedNews := []dto.NewsResponseDTO{
			{
				ID:         testNewsID,
				CategoryID: testNewsCategoryID,
				AuthorID:   testNewsAuthorID,
				Title:      "Breaking News",
				Content:    "This is the news content",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}

		// Mock expectations
		mockNewsUseCase.On("GetAll", mock.Anything).Return(expectedNews, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/news", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.GET("/news", newsRouter.GetAll)

		// Mock expectations
		mockNewsUseCase.On("GetAll", mock.Anything).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodGet, "/news", http.NoBody)
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

		mockNewsUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestNewsRoutes_GetByID(t *testing.T) {
	t.Run("success - get news by id", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.GET("/news/:id", newsRouter.GetByID)

		now := time.Now()
		expectedNews := &dto.NewsResponseDTO{
			ID:         testNewsID,
			CategoryID: testNewsCategoryID,
			AuthorID:   testNewsAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Mock expectations
		mockNewsUseCase.On("GetByID", mock.Anything, testNewsID).Return(expectedNews, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/news/"+testNewsID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - news not found", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.GET("/news/:id", newsRouter.GetByID)

		// Mock expectations
		mockNewsUseCase.On("GetByID", mock.Anything, "non-existent-id").Return(nil, apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/news/non-existent-id", http.NoBody)
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
		assert.Equal(t, "News not found", meta["message"])

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.GET("/news/:id", newsRouter.GetByID)

		// Mock expectations
		mockNewsUseCase.On("GetByID", mock.Anything, testNewsID).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodGet, "/news/"+testNewsID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockNewsUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestNewsRoutes_Create(t *testing.T) {
	t.Run("success - create news", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		// Simulate authenticated user middleware setting user_id
		router.POST("/news", func(c *gin.Context) {
			c.Set("user_id", testNewsAuthorID)
			newsRouter.Create(c)
		})

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Breaking News",
			"content":     "This is the news content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		now := time.Now()
		expectedNews := &dto.NewsResponseDTO{
			ID:         testNewsID,
			CategoryID: testNewsCategoryID,
			AuthorID:   testNewsAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Mock expectations
		mockNewsUseCase.On("Create", mock.Anything, testNewsAuthorID, &dto.CreateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Breaking News",
			Content:    "This is the news content",
		}).Return(expectedNews, nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(bodyBytes))
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

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.POST("/news", func(c *gin.Context) {
			c.Set("user_id", testNewsAuthorID)
			newsRouter.Create(c)
		})

		// Missing required fields
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(bodyBytes))
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

		mockNewsUseCase.AssertNotCalled(t, "Create")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - user not authenticated", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.POST("/news", newsRouter.Create)

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Breaking News",
			"content":     "This is the news content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(bodyBytes))
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

		mockNewsUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("error - invalid user id format", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.POST("/news", func(c *gin.Context) {
			c.Set("user_id", 12345) // Set as int instead of string
			newsRouter.Create(c)
		})

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Breaking News",
			"content":     "This is the news content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(bodyBytes))
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

		mockNewsUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.POST("/news", func(c *gin.Context) {
			c.Set("user_id", testNewsAuthorID)
			newsRouter.Create(c)
		})

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Breaking News",
			"content":     "This is the news content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockNewsUseCase.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil, apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockNewsUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestNewsRoutes_Update(t *testing.T) {
	t.Run("success - update news", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.PUT("/news/:id", newsRouter.Update)

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Updated News",
			"content":     "Updated content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockNewsUseCase.On("Update", mock.Anything, testNewsID, &dto.UpdateNewsRequestDTO{
			CategoryID: testNewsCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/news/"+testNewsID, bytes.NewBuffer(bodyBytes))
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

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.PUT("/news/:id", newsRouter.Update)

		// Missing required fields
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/news/"+testNewsID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockNewsUseCase.AssertNotCalled(t, "Update")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - news not found", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.PUT("/news/:id", newsRouter.Update)

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Updated News",
			"content":     "Updated content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockNewsUseCase.On("Update", mock.Anything, "non-existent-id", mock.Anything).Return(apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/news/non-existent-id", bytes.NewBuffer(bodyBytes))
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
		assert.Equal(t, "News not found", meta["message"])

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.PUT("/news/:id", newsRouter.Update)

		requestBody := map[string]string{
			"category_id": testNewsCategoryID,
			"title":       "Updated News",
			"content":     "Updated content",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockNewsUseCase.On("Update", mock.Anything, testNewsID, mock.Anything).Return(apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/news/"+testNewsID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockNewsUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestNewsRoutes_Delete(t *testing.T) {
	t.Run("success - delete news", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.DELETE("/news/:id", newsRouter.Delete)

		// Mock expectations
		mockNewsUseCase.On("Delete", mock.Anything, testNewsID).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/news/"+testNewsID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])
		assert.NotNil(t, response["data"])

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - news not found", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.DELETE("/news/:id", newsRouter.Delete)

		// Mock expectations
		mockNewsUseCase.On("Delete", mock.Anything, "non-existent-id").Return(apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/news/non-existent-id", http.NoBody)
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
		assert.Equal(t, "News not found", meta["message"])

		mockNewsUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockNewsUseCase := new(MockNewsUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		newsRouter := &newsRoutes{
			news: mockNewsUseCase,
			log:  mockLogger,
		}

		router.DELETE("/news/:id", newsRouter.Delete)

		// Mock expectations
		mockNewsUseCase.On("Delete", mock.Anything, testNewsID).Return(apperror.ErrDatabaseConnection)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/news/"+testNewsID, http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockNewsUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
