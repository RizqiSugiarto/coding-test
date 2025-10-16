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

type customPageRoutes struct {
	customPage usecase.CustomPage
	log        logger.Interface
}

func newCustomPageRoutes(handler *gin.RouterGroup, customPage usecase.CustomPage, log logger.Interface, authMiddleware gin.HandlerFunc) {
	customPageRouter := customPageRoutes{customPage, log}

	h := handler.Group("pages")
	{
		// Public endpoints - anyone can read custom pages
		h.GET("", customPageRouter.GetAll)
		h.GET("/:id", customPageRouter.GetByID)

		// Protected endpoints - only authenticated users
		h.POST("", authMiddleware, customPageRouter.Create)
		h.PUT("/:id", authMiddleware, customPageRouter.Update)
		h.DELETE("/:id", authMiddleware, customPageRouter.Delete)
	}
}

// @Summary Get all custom pages
// @Description Retrieve a list of all custom pages
// @Tags CustomPages
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "List of custom pages"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /pages [get]
func (cp *customPageRoutes) GetAll(ctx *gin.Context) {
	pageList, err := cp.customPage.GetAll(ctx)
	if err != nil {
		cp.log.Error(err, "CustomPageController - GetAll - cp.customPage.GetAll")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"pages": pageList,
	})
}

// @Summary Get custom page by ID
// @Description Retrieve a single custom page by its ID
// @Tags CustomPages
// @Accept json
// @Produce json
// @Param id path string true "Page ID"
// @Success 200 {object} response.Response "Page detail"
// @Failure 404 {object} response.ErrorResponse "Page not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /pages/{id} [get]
func (cp *customPageRoutes) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	page, err := cp.customPage.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "Page not found")

			return
		}

		cp.log.Error(err, "CustomPageController - GetByID - cp.customPage.GetByID")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"page": page,
	})
}

// @Summary Create a new custom page
// @Description Create a new custom page (requires authentication)
// @Tags CustomPages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CustomPage true "Page information"
// @Success 201 {object} response.Response "Page created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /pages [post]
func (cp *customPageRoutes) Create(ctx *gin.Context) {
	var req request.CustomPage

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		cp.log.Error(err, "CustomPageController - Create - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.SendError(ctx, http.StatusUnauthorized, "User not authenticated")

		return
	}

	authorID, ok := userID.(string)
	if !ok {
		response.SendError(ctx, http.StatusInternalServerError, "Invalid user ID format")

		return
	}

	// Create custom page
	page, err := cp.customPage.Create(ctx, authorID, &dto.CreateCustomPageRequestDTO{
		CustomURL: req.CustomURL,
		Content:   req.Content,
	})
	if err != nil {
		cp.log.Error(err, "CustomPageController - Create - cp.customPage.Create")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusCreated, gin.H{
		"page": page,
	})
}

// @Summary Update a custom page
// @Description Update an existing custom page (requires authentication)
// @Tags CustomPages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Page ID"
// @Param request body request.UpdateCustomPage true "Updated page information"
// @Success 200 {object} response.Response "Page updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Page not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /pages/{id} [put]
func (cp *customPageRoutes) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req request.UpdateCustomPage

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		cp.log.Error(err, "CustomPageController - Update - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Update custom page
	err := cp.customPage.Update(ctx, id, &dto.UpdateCustomPageRequestDTO{
		CustomURL: req.CustomURL,
		Content:   req.Content,
	})
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "Page not found")

			return
		}

		cp.log.Error(err, "CustomPageController - Update - cp.customPage.Update")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"message": "Page updated successfully",
	})
}

// @Summary Delete a custom page
// @Description Delete a custom page by ID (requires authentication)
// @Tags CustomPages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Page ID"
// @Success 200 {object} response.Response "Page deleted successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Page not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /pages/{id} [delete]
func (cp *customPageRoutes) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Delete custom page
	err := cp.customPage.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "Page not found")

			return
		}

		cp.log.Error(err, "CustomPageController - Delete - cp.customPage.Delete")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"message": "Page deleted successfully",
	})
}
