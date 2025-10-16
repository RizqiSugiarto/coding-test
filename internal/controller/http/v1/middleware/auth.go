package middleware

import (
	"net/http"
	"strings"

	"github.com/RizqiSugiarto/coding-test/internal/controller/http/v1/response"
	"github.com/RizqiSugiarto/coding-test/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userIDKey           = "user_id"
)

// AuthMiddleware creates a middleware for JWT authentication.
func AuthMiddleware(jwtManager jwt.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader(authorizationHeader)
		if authHeader == "" {
			response.SendError(ctx, http.StatusUnauthorized, "Authorization header is required")
			ctx.Abort()

			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.SendError(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			ctx.Abort()

			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if token == "" {
			response.SendError(ctx, http.StatusUnauthorized, "Token is required")
			ctx.Abort()

			return
		}

		// Validate token
		claims, err := jwtManager.ParseAndValidateAccessToken(token)
		if err != nil {
			response.SendError(ctx, http.StatusUnauthorized, "Invalid or expired token")
			ctx.Abort()

			return
		}

		// Extract user_id from claims
		userID, ok := claims[userIDKey].(string)
		if !ok {
			response.SendError(ctx, http.StatusUnauthorized, "Invalid token claims")
			ctx.Abort()

			return
		}

		// Set user_id in context
		ctx.Set(userIDKey, userID)

		ctx.Next()
	}
}
