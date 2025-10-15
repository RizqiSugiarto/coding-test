package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUseCase is a mock implementation of usecase.Auth
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Login(ctx context.Context, req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*dto.LoginResponseDTO)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

// MockLogger is a mock implementation of logger.Interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Info(message string, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Warn(message string, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Error(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	return router
}

func TestAuthRoutes_Login(t *testing.T) {
	t.Run("success - valid credentials", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		requestBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		expectedResponse := &dto.LoginResponseDTO{
			AccessToken:  "access.token.here",
			RefreshToken: "refresh.token.here",
		}

		// Mock expectations
		mockAuthUseCase.On("Login", mock.Anything, dto.LoginRequestDTO{
			UserName: "testuser",
			Password: "password123",
		}).Return(expectedResponse, nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
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

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(200), meta["code"])
		assert.Equal(t, "OK", meta["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok, "data should be a map")

		token, ok := data["token"].(map[string]interface{})
		assert.True(t, ok, "token should be a map")
		assert.Equal(t, expectedResponse.AccessToken, token["access_token"])
		assert.Equal(t, expectedResponse.RefreshToken, token["refresh_token"])

		mockAuthUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid request payload (malformed JSON)", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		// Malformed JSON
		bodyBytes := []byte(`{"username": "testuser", "password":}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(400), meta["code"])
		assert.Equal(t, "Invalid request payload", meta["message"])

		mockAuthUseCase.AssertNotCalled(t, "Login")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - invalid request payload (missing required fields)", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		// Missing password field
		requestBody := map[string]string{
			"username": "testuser",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(400), meta["code"])
		assert.Equal(t, "Invalid request payload", meta["message"])

		mockAuthUseCase.AssertNotCalled(t, "Login")
		mockLogger.AssertExpectations(t)
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		requestBody := map[string]string{
			"username": "nonexistentuser",
			"password": "password123",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockAuthUseCase.On("Login", mock.Anything, dto.LoginRequestDTO{
			UserName: "nonexistentuser",
			Password: "password123",
		}).Return(nil, apperror.ErrNotFound)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(404), meta["code"])
		assert.Equal(t, "User not found", meta["message"])

		mockAuthUseCase.AssertExpectations(t)
	})

	t.Run("error - invalid credentials", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		requestBody := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockAuthUseCase.On("Login", mock.Anything, dto.LoginRequestDTO{
			UserName: "testuser",
			Password: "wrongpassword",
		}).Return(nil, apperror.ErrInvalidCredentials)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(401), meta["code"])
		assert.Equal(t, "Invalid username or password", meta["message"])

		mockAuthUseCase.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		requestBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Mock expectations
		mockAuthUseCase.On("Login", mock.Anything, dto.LoginRequestDTO{
			UserName: "testuser",
			Password: "password123",
		}).Return(nil, apperror.ErrDatabaseConnection)

		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(500), meta["code"])
		assert.Equal(t, "Internal server error", meta["message"])

		mockAuthUseCase.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("success - login with empty request body", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		bodyBytes := []byte(`{}`)

		// Mock expectations
		mockLogger.On("Error", mock.Anything, mock.Anything).Return()

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["meta"])

		meta, ok := response["meta"].(map[string]interface{})
		assert.True(t, ok, "meta should be a map")
		assert.Equal(t, float64(400), meta["code"])
		assert.Equal(t, "Invalid request payload", meta["message"])

		mockAuthUseCase.AssertNotCalled(t, "Login")
		mockLogger.AssertExpectations(t)
	})

	t.Run("success - login with special characters in credentials", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		authRouter := &authRoutes{
			auth: mockAuthUseCase,
			log:  mockLogger,
		}

		router.POST("/auth/login", authRouter.Login)

		requestBody := map[string]string{
			"username": "test.user+special@example.com",
			"password": "P@ssw0rd!#$",
		}
		bodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		expectedResponse := &dto.LoginResponseDTO{
			AccessToken:  "access.token.special",
			RefreshToken: "refresh.token.special",
		}

		// Mock expectations
		mockAuthUseCase.On("Login", mock.Anything, dto.LoginRequestDTO{
			UserName: "test.user+special@example.com",
			Password: "P@ssw0rd!#$",
		}).Return(expectedResponse, nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
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

		mockAuthUseCase.AssertExpectations(t)
	})
}

func TestNewAuthRoutes(t *testing.T) {
	t.Run("success - create auth routes", func(t *testing.T) {
		// Arrange
		mockAuthUseCase := new(MockAuthUseCase)
		mockLogger := new(MockLogger)

		router := setupTestRouter()
		handler := router.Group("/api/v1")

		// Mock logger since the endpoint will be called and may fail with bad request
		mockLogger.On("Error", mock.Anything, mock.Anything).Return().Maybe()

		// Act
		newAuthRoutes(handler, mockAuthUseCase, mockLogger)

		// Assert - verify route is registered
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return 400 (bad request) instead of 404 (not found)
		// This proves the route exists
		assert.NotEqual(t, http.StatusNotFound, w.Code)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
