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

type newsRoutes struct {
	news usecase.News
	log  logger.Interface
}

func newNewsRoutes(handler *gin.RouterGroup, news usecase.News, log logger.Interface, authMiddleware gin.HandlerFunc) {
	newsRouter := newsRoutes{news, log}

	h := handler.Group("news")
	{
		// Public endpoints - anyone can read news
		h.GET("", newsRouter.GetAll)
		h.GET("/:id", newsRouter.GetByID)

		// Protected endpoints - only authenticated users
		h.POST("", authMiddleware, newsRouter.Create)
		h.PUT("/:id", authMiddleware, newsRouter.Update)
		h.DELETE("/:id", authMiddleware, newsRouter.Delete)
	}
}

// @Summary Get all news
// @Description Retrieve a list of all news articles
// @Tags News
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "List of news"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news [get]
func (n *newsRoutes) GetAll(ctx *gin.Context) {
	newsList, err := n.news.GetAll(ctx)
	if err != nil {
		n.log.Error(err, "NewsController - GetAll - n.news.GetAll")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"news": newsList,
	})
}

// @Summary Get news by ID
// @Description Retrieve a single news article by its ID
// @Tags News
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} response.Response "News detail"
// @Failure 404 {object} response.ErrorResponse "News not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news/{id} [get]
func (n *newsRoutes) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	news, err := n.news.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "News not found")

			return
		}

		n.log.Error(err, "NewsController - GetByID - n.news.GetByID")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"news": news,
	})
}

// @Summary Create a new news article
// @Description Create a new news article (requires authentication)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.News true "News information"
// @Success 201 {object} response.Response "News created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news [post]
func (n *newsRoutes) Create(ctx *gin.Context) {
	var req request.News

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		n.log.Error(err, "NewsController - Create - ctx.ShouldBindJSON")
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

	// Create news
	news, err := n.news.Create(ctx, authorID, &dto.CreateNewsRequestDTO{
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
	})
	if err != nil {
		n.log.Error(err, "NewsController - Create - n.news.Create")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusCreated, gin.H{
		"news": news,
	})
}

// @Summary Update a news article
// @Description Update an existing news article (requires authentication)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "News ID"
// @Param request body request.UpdateNews true "Updated news information"
// @Success 200 {object} response.Response "News updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "News not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news/{id} [put]
func (n *newsRoutes) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req request.UpdateNews

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		n.log.Error(err, "NewsController - Update - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Update news
	err := n.news.Update(ctx, id, &dto.UpdateNewsRequestDTO{
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
	})
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "News not found")

			return
		}

		n.log.Error(err, "NewsController - Update - n.news.Update")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"message": "News updated successfully",
	})
}

// @Summary Delete a news article
// @Description Delete a news article by ID (requires authentication)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "News ID"
// @Success 200 {object} response.Response "News deleted successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "News not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news/{id} [delete]
func (n *newsRoutes) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Delete news
	err := n.news.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.SendError(ctx, http.StatusNotFound, "News not found")

			return
		}

		n.log.Error(err, "NewsController - Delete - n.news.Delete")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusOK, gin.H{
		"message": "News deleted successfully",
	})
}
