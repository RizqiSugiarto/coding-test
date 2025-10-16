package usecase

import (
	"context"
	"testing"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const (
	testPassword          = "password123"
	testAccessToken       = "access.token.here"
	testValidRefreshToken = "valid.refresh.token"
	testUserID            = "user-123"
	testNewRefreshToken   = "new.refresh.token"
)

// MockUserRepo is a mock implementation of repository.UserRepo.
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

// MockJWTManager is a mock implementation of jwt.Manager.
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) GenerateAccessToken(userID string) (string, error) {
	args := m.Called(userID)

	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) GenerateRefreshToken(userID string) (string, error) {
	args := m.Called(userID)

	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) ParseAndValidateAccessToken(tokenStr string) (jwt.MapClaims, error) {
	args := m.Called(tokenStr)

	if claims, ok := args.Get(0).(jwt.MapClaims); ok {
		return claims, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockJWTManager) ParseAndValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	args := m.Called(tokenStr)

	if claims, ok := args.Get(0).(jwt.MapClaims); ok {
		return claims, args.Error(1)
	}

	return nil, args.Error(1)
}

// Helper function to generate bcrypt hash for testing.
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) // This should never happen in tests
	}

	return string(hash)
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Run("success - valid credentials", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		password := testPassword
		hashedPassword := hashPassword(password)

		expectedUser := &entity.User{
			ID:       "user-123",
			Username: "testuser",
			Password: hashedPassword,
		}

		loginReq := dto.LoginRequestDTO{
			UserName: "testuser",
			Password: password,
		}

		expectedAccessToken := testAccessToken
		expectedRefreshToken := "refresh.token.here"

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(expectedUser, nil)
		mockJWTManager.On("GenerateAccessToken", expectedUser.ID).Return(expectedAccessToken, nil)
		mockJWTManager.On("GenerateRefreshToken", expectedUser.ID).Return(expectedRefreshToken, nil)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedAccessToken, resp.AccessToken)
		assert.Equal(t, expectedRefreshToken, resp.RefreshToken)
		mockUserRepo.AssertExpectations(t)
		mockJWTManager.AssertExpectations(t)
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		loginReq := dto.LoginRequestDTO{
			UserName: "nonexistentuser",
			Password: testPassword,
		}

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(nil, apperror.ErrNotFound)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTService.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - repository error", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		loginReq := dto.LoginRequestDTO{
			UserName: "testuser",
			Password: testPassword,
		}

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(nil, apperror.ErrDatabaseConnection)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTService.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - invalid password", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		correctPassword := "correctpassword"
		wrongPassword := "wrongpassword"
		hashedPassword := hashPassword(correctPassword)

		expectedUser := &entity.User{
			ID:       "user-123",
			Username: "testuser",
			Password: hashedPassword,
		}

		loginReq := dto.LoginRequestDTO{
			UserName: "testuser",
			Password: wrongPassword,
		}

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(expectedUser, nil)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrInvalidCredentials, err)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTService.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - access token generation fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		password := testPassword
		hashedPassword := hashPassword(password)

		expectedUser := &entity.User{
			ID:       "user-123",
			Username: "testuser",
			Password: hashedPassword,
		}

		loginReq := dto.LoginRequestDTO{
			UserName: "testuser",
			Password: password,
		}

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(expectedUser, nil)
		mockJWTService.On("GenerateAccessToken", expectedUser.ID).Return("", apperror.ErrGenerateAccessToken)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrGenerateAccessToken, err)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - refresh token generation fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		password := testPassword
		hashedPassword := hashPassword(password)

		expectedUser := &entity.User{
			ID:       "user-123",
			Username: "testuser",
			Password: hashedPassword,
		}

		loginReq := dto.LoginRequestDTO{
			UserName: "testuser",
			Password: password,
		}

		expectedAccessToken := testAccessToken

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(expectedUser, nil)
		mockJWTService.On("GenerateAccessToken", expectedUser.ID).Return(expectedAccessToken, nil)
		mockJWTService.On("GenerateRefreshToken", expectedUser.ID).Return("", apperror.ErrGenerateRefreshToken)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrGenerateRefreshToken, err)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	t.Run("success - login with special characters in username", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		password := "P@ssw0rd!#$"
		hashedPassword := hashPassword(password)

		expectedUser := &entity.User{
			ID:       "user-456",
			Username: "test.user+special@example.com",
			Password: hashedPassword,
		}

		loginReq := dto.LoginRequestDTO{
			UserName: "test.user+special@example.com",
			Password: password,
		}

		expectedAccessToken := "access.token.special"
		expectedRefreshToken := "refresh.token.special"

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(expectedUser, nil)
		mockJWTService.On("GenerateAccessToken", expectedUser.ID).Return(expectedAccessToken, nil)
		mockJWTService.On("GenerateRefreshToken", expectedUser.ID).Return(expectedRefreshToken, nil)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedAccessToken, resp.AccessToken)
		assert.Equal(t, expectedRefreshToken, resp.RefreshToken)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	t.Run("error - empty username", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		loginReq := dto.LoginRequestDTO{
			UserName: "",
			Password: testPassword,
		}

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(nil, apperror.ErrNotFound)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - empty password", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		ctx := context.Background()
		password := testPassword
		hashedPassword := hashPassword(password)

		expectedUser := &entity.User{
			ID:       "user-123",
			Username: "testuser",
			Password: hashedPassword,
		}

		loginReq := dto.LoginRequestDTO{
			UserName: "testuser",
			Password: "", // Empty password
		}

		// Mock expectations
		mockUserRepo.On("GetByUsername", ctx, loginReq.UserName).Return(expectedUser, nil)

		// Act
		resp, err := authUseCase.Login(ctx, loginReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrInvalidCredentials, err)
		mockUserRepo.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTService.AssertNotCalled(t, "GenerateRefreshToken")
	})
}

func TestAuthUseCase_Refresh(t *testing.T) {
	t.Run("success - valid refresh token", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := testValidRefreshToken
		userID := testUserID

		claims := jwt.MapClaims{
			"user_id": userID,
		}

		expectedAccessToken := testAccessToken
		expectedRefreshToken := testNewRefreshToken

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(claims, nil)
		mockJWTManager.On("GenerateAccessToken", userID).Return(expectedAccessToken, nil)
		mockJWTManager.On("GenerateRefreshToken", userID).Return(expectedRefreshToken, nil)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedAccessToken, resp.AccessToken)
		assert.Equal(t, expectedRefreshToken, resp.RefreshToken)
		mockJWTManager.AssertExpectations(t)
	})

	t.Run("error - invalid refresh token", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := "invalid.refresh.token"

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(nil, apperror.ErrInvalidToken)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrInvalidToken, err)
		mockJWTManager.AssertExpectations(t)
		mockJWTManager.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTManager.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - expired refresh token", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := "expired.refresh.token"

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(nil, apperror.ErrInvalidToken)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrInvalidToken, err)
		mockJWTManager.AssertExpectations(t)
		mockJWTManager.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTManager.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - access token generation fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := testValidRefreshToken
		userID := testUserID

		claims := jwt.MapClaims{
			"user_id": userID,
		}

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(claims, nil)
		mockJWTManager.On("GenerateAccessToken", userID).Return("", apperror.ErrGenerateAccessToken)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrGenerateAccessToken, err)
		mockJWTManager.AssertExpectations(t)
		mockJWTManager.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("error - refresh token generation fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := testValidRefreshToken
		userID := testUserID

		claims := jwt.MapClaims{
			"user_id": userID,
		}

		expectedAccessToken := testAccessToken

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(claims, nil)
		mockJWTManager.On("GenerateAccessToken", userID).Return(expectedAccessToken, nil)
		mockJWTManager.On("GenerateRefreshToken", userID).Return("", apperror.ErrGenerateRefreshToken)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrGenerateRefreshToken, err)
		mockJWTManager.AssertExpectations(t)
	})

	t.Run("error - invalid token type", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := "access.token.instead.of.refresh"

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(nil, apperror.ErrInvalidTokenType)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, apperror.ErrInvalidTokenType, err)
		mockJWTManager.AssertExpectations(t)
		mockJWTManager.AssertNotCalled(t, "GenerateAccessToken")
		mockJWTManager.AssertNotCalled(t, "GenerateRefreshToken")
	})

	t.Run("success - refresh with different user ID format", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepo)
		mockJWTManager := new(MockJWTManager)
		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTManager)

		ctx := context.Background()
		refreshToken := testValidRefreshToken
		userID := "550e8400-e29b-41d4-a716-446655440000" // UUID format

		claims := jwt.MapClaims{
			"user_id": userID,
		}

		expectedAccessToken := "new.access.token"
		expectedRefreshToken := testNewRefreshToken

		// Mock expectations
		mockJWTManager.On("ParseAndValidateRefreshToken", refreshToken).Return(claims, nil)
		mockJWTManager.On("GenerateAccessToken", userID).Return(expectedAccessToken, nil)
		mockJWTManager.On("GenerateRefreshToken", userID).Return(expectedRefreshToken, nil)

		// Act
		resp, err := authUseCase.Refresh(ctx, refreshToken)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedAccessToken, resp.AccessToken)
		assert.Equal(t, expectedRefreshToken, resp.RefreshToken)
		mockJWTManager.AssertExpectations(t)
	})
}

// TestNewAuthUseCase tests the constructor.
func TestNewAuthUseCase(t *testing.T) {
	t.Run("success - create new auth usecase", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockJWTService := new(MockJWTManager)

		authUseCase := NewAuthUseCase(mockUserRepo, mockJWTService)

		assert.NotNil(t, authUseCase)
		assert.NotNil(t, authUseCase.userRepo)
		assert.NotNil(t, authUseCase.jwtManager)
	})
}
