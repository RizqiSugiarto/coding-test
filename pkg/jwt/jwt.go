package jwt

import (
	"time"

	"github.com/RizqiSugiarto/coding-test/config"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
)

// Manager defines the interface for JWT token operations.
type Manager interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ParseAndValidateAccessToken(tokenStr string) (jwt.MapClaims, error)
	ParseAndValidateRefreshToken(tokenStr string) (jwt.MapClaims, error)
}

type manager struct {
	config *config.JWT
}

// NewJWTManager creates a new JWT manager instance.
func NewJWTManager(cfg *config.JWT) Manager {
	return &manager{config: cfg}
}

func (j *manager) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.config.AccessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.config.AccessTokenSecretKey))
}

func (j *manager) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(j.config.RefreshTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.config.RefreshTokenSecretKey))
}

func (j *manager) ParseAndValidateAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(j.config.AccessTokenSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.ErrInvalidTokenClaims
	}

	return claims, nil
}

func (j *manager) ParseAndValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(j.config.RefreshTokenSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, apperror.ErrInvalidTokenType
	}

	return claims, nil
}
