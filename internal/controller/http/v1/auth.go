package v1

import (
	"errors"
	"net/http"

	"github.com/RizqiSugiarto/coding-test/internal/controller/http/v1/request"
	"github.com/RizqiSugiarto/coding-test/internal/controller/http/v1/response"
	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/usecase"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/logger"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	auth usecase.Auth
	log  logger.Interface
}

func newAuthRoutes(handler *gin.RouterGroup, auth usecase.Auth, log logger.Interface) {
	authRouter := authRoutes{auth, log}

	h := handler.Group("auth")
	{
		h.POST("/login", authRouter.Login)
	}
}

// @Summary User login
// @Description Authenticate user with username and password, returning a JWT token if valid.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.Auth true "Login credentials"
// @Success 200 {object} response.LoginSuccessResponse "JWT token"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Invalid username or password"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (a *authRoutes) Login(ctx *gin.Context) {
	var req request.Auth

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		a.log.Error(err, "AuthController - Login - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Perform login
	token, err := a.auth.Login(ctx, dto.LoginRequestDTO{
		UserName: req.Username,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrInvalidCredentials):
			response.SendError(ctx, http.StatusUnauthorized, "Invalid username or password")
		case errors.Is(err, apperror.ErrNotFound):
			response.SendError(ctx, http.StatusNotFound, "User not found")
		default:
			a.log.Error(err, "AuthController - Login - a.auth.Login")
			response.SendError(ctx, http.StatusInternalServerError, "Internal server error")
		}

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"token": token,
	})
}
