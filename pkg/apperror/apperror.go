package apperror

import "errors"

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrNotFound             = errors.New("resource not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidTokenClaims   = errors.New("invalid token claims")
	ErrInvalidTokenType     = errors.New("invalid token type")
	ErrDatabaseConnection   = errors.New("database connection failed")
	ErrGenerateAccessToken  = errors.New("failed to generate access token")
	ErrGenerateRefreshToken = errors.New("failed to generate refresh token")
	ErrDuplicateKey         = errors.New("duplicate key value violates unique constraint")
)
