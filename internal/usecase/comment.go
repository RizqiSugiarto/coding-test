package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/internal/repository"
)

type CommentUseCase struct {
	commentRepo repository.CommentRepo
}

func NewCommentUseCase(commentRepo repository.CommentRepo) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
	}
}

func (co *CommentUseCase) Create(ctx context.Context, req *dto.CreateCommentRequestDTO) error {
	comment := &entity.Comment{
		Name:    req.Name,
		NewsID:  req.NewsID,
		Comment: req.Comment,
	}

	err := co.commentRepo.Create(ctx, comment)
	if err != nil {
		return err
	}

	return nil
}
