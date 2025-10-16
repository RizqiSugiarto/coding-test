package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testCommentNewsIDRoute = "550e8400-e29b-41d4-a716-446655440000"

var errCommentDatabase = errors.New("database error")

// MockCommentUseCase is a mock implementation of usecase.Comment.
type MockCommentUseCase struct {
	mock.Mock
}

func (m *MockCommentUseCase) Create(ctx context.Context, req *dto.CreateCommentRequestDTO) error {
	args := m.Called(ctx, req)

	return args.Error(0)
}

func TestCommentRoutes_Create(t *testing.T) {
	t.Run("success - create comment", func(t *testing.T) {
		// Arrange
		mockCommentUseCase := new(MockCommentUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		commentRouter := &commentRoutes{
			comment: mockCommentUseCase,
			log:     mockLogger,
		}

		router.POST("/news/:id/comments", commentRouter.Create)

		requestBody := map[string]string{
			"name":    "John Doe",
			"comment": "This is a great article!",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCommentUseCase.On("Create", mock.Anything, &dto.CreateCommentRequestDTO{
			Name:    "John Doe",
			Comment: "This is a great article!",
			NewsID:  testCommentNewsIDRoute,
		}).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news/"+testCommentNewsIDRoute+"/comments", bytes.NewBuffer(bodyBytes))
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

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "Comment created successfully", data["message"])

		mockCommentUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload (missing required fields)", func(t *testing.T) {
		// Arrange
		mockCommentUseCase := new(MockCommentUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		commentRouter := &commentRoutes{
			comment: mockCommentUseCase,
			log:     mockLogger,
		}

		router.POST("/news/:id/comments", commentRouter.Create)

		// Missing required fields
		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news/"+testCommentNewsIDRoute+"/comments", bytes.NewBuffer(bodyBytes))
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

		mockCommentUseCase.AssertNotCalled(t, "Create")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - invalid request payload (malformed JSON)", func(t *testing.T) {
		// Arrange
		mockCommentUseCase := new(MockCommentUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		commentRouter := &commentRoutes{
			comment: mockCommentUseCase,
			log:     mockLogger,
		}

		router.POST("/news/:id/comments", commentRouter.Create)

		// Malformed JSON
		bodyBytes := []byte(`{"name": "John Doe", "comment":}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news/"+testCommentNewsIDRoute+"/comments", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockCommentUseCase.AssertNotCalled(t, "Create")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockCommentUseCase := new(MockCommentUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		commentRouter := &commentRoutes{
			comment: mockCommentUseCase,
			log:     mockLogger,
		}

		router.POST("/news/:id/comments", commentRouter.Create)

		requestBody := map[string]string{
			"name":    "John Doe",
			"comment": "This is a great article!",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCommentUseCase.On("Create", mock.Anything, mock.Anything).Return(errCommentDatabase)
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news/"+testCommentNewsIDRoute+"/comments", bytes.NewBuffer(bodyBytes))
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
		assert.Equal(t, "Internal server error", meta["message"])

		mockCommentUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("success - create comment with special characters", func(t *testing.T) {
		// Arrange
		mockCommentUseCase := new(MockCommentUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		commentRouter := &commentRoutes{
			comment: mockCommentUseCase,
			log:     mockLogger,
		}

		router.POST("/news/:id/comments", commentRouter.Create)

		requestBody := map[string]string{
			"name":    "Test User <script>",
			"comment": "Comment with special chars: !@#$%^&*()",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockCommentUseCase.On("Create", mock.Anything, &dto.CreateCommentRequestDTO{
			Name:    "Test User <script>",
			Comment: "Comment with special chars: !@#$%^&*()",
			NewsID:  testCommentNewsIDRoute,
		}).Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/news/"+testCommentNewsIDRoute+"/comments", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		mockCommentUseCase.AssertExpectations(t)
	})
}
