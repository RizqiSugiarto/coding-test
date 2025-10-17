package v1

import (
	"net/http"

	"github.com/RizqiSugiarto/coding-test/internal/controller/http/v1/request"
	"github.com/RizqiSugiarto/coding-test/internal/controller/http/v1/response"
	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/usecase"
	"github.com/RizqiSugiarto/coding-test/pkg/logger"
	"github.com/gin-gonic/gin"
)

type commentRoutes struct {
	comment usecase.Comment
	log     logger.Interface
}

func newCommentRoutes(handler *gin.RouterGroup, comment usecase.Comment, log logger.Interface) {
	commentRouter := commentRoutes{comment, log}

	h := handler.Group("news")
	{
		// Public endpoint - anyone can post comments
		h.POST("/:id/comments", commentRouter.Create)
	}
}

// @Summary Create a comment on a news article
// @Description Create a new comment on a specific news article
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Param request body request.Comment true "Comment information"
// @Success 201 {object} response.Response "Comment created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request payload"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news/{id}/comments [post]
func (co *commentRoutes) Create(ctx *gin.Context) {
	newsID := ctx.Param("id")

	var req request.Comment

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		co.log.Error(err, "CommentController - Create - ctx.ShouldBindJSON")
		response.SendError(ctx, http.StatusBadRequest, "Invalid request payload")

		return
	}

	// Create comment
	err := co.comment.Create(ctx, &dto.CreateCommentRequestDTO{
		Name:    req.Name,
		Comment: req.Comment,
		NewsID:  newsID,
	})
	if err != nil {
		co.log.Error(err, "CommentController - Create - co.comment.Create")
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error")

		return
	}

	// Success response
	response.SendSuccess(ctx, http.StatusCreated, gin.H{
		"message": "Comment created successfully",
	})
}
