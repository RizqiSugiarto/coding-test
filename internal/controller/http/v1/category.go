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

type categoryRoutes struct {
	category usecase.Category
	log      logger.Interface
}

func newCategoryRoutes(handler *gin.RouterGroup, category usecase.Category, log logger.Interface, authMiddleware gin.HandlerFunc) {
	categoryRouter := categoryRoutes{category, log}

	h := handler.Group("categories")
	{
		// Public endpoints - anyone can read categories
		h.GET("", categoryRouter.GetAll)
		h.GET("/:id", categoryRouter.GetByID)

		// Protected endpoints - only authenticated users
		h.POST("", authMiddleware, categoryRouter.Create)
		h.PUT("/:id", authMiddleware, categoryRouter.Update)
		h.DELETE("/:id", authMiddleware, categoryRouter.Delete)
	}
}

// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "List of categories"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories [get]
func (c *categoryRoutes) GetAll(ctx *gin.Context) {
	categories, err := c.category.GetAll(ctx)
	if err != nil {
		c.log.Error(err, "CategoryController - GetAll - c.category.GetAll")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"categories": categories,
	})
}

// @Summary Get category by ID
// @Description Retrieve a single category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response "Category detail"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories/{id} [get]
func (c *categoryRoutes) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	category, err := c.category.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "Category not found")

			return
		}

		c.log.Error(err, "CategoryController - GetByID - c.category.GetByID")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"category": category,
	})
}

// @Summary Create a new category
// @Description Create a new category (requires authentication)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.Category true "Category information"
// @Success 201 {object} response.Response "Category created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories [post]
func (c *categoryRoutes) Create(ctx *gin.Context) {
	var req request.Category

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.log.Error(err, "CategoryController - Create - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Create category
	category, err := c.category.Create(ctx, &dto.CreateCategoryRequestDTO{
		Name: req.Name,
	})
	if err != nil {
		c.log.Error(err, "CategoryController - Create - c.category.Create")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusCreated, gin.H{
		"category": category,
	})
}

// @Summary Update a category
// @Description Update an existing category (requires authentication)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body request.UpdateCategory true "Updated category information"
// @Success 200 {object} response.Response "Category updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories/{id} [put]
func (c *categoryRoutes) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req request.UpdateCategory

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.log.Error(err, "CategoryController - Update - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Update category
	err := c.category.Update(ctx, id, &dto.UpdateCategoryRequestDTO{
		Name: req.Name,
	})
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "Category not found")

			return
		}

		c.log.Error(err, "CategoryController - Update - c.category.Update")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"message": "Category updated successfully",
	})
}

// @Summary Delete a category
// @Description Delete a category by ID (requires authentication)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response "Category deleted successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories/{id} [delete]
func (c *categoryRoutes) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Delete category
	err := c.category.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "Category not found")

			return
		}

		c.log.Error(err, "CategoryController - Delete - c.category.Delete")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
