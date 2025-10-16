package postgres

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/postgres"
)

type CommentRepo struct {
	*postgres.Postgres
}

func NewPostgresCommentRepo(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

func (c *CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	// Check if news exists
	var exists bool

	checkSQL, _, err := c.Builder.Select("EXISTS(SELECT 1 FROM news WHERE id = ?)").ToSql()
	if err != nil {
		return err
	}

	err = c.DB.QueryRowContext(ctx, checkSQL, comment.NewsID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return apperror.ErrNotFound
	}

	// Insert comment
	sql, args, err := c.Builder.Insert("comments").
		Columns("name, news_id, comment").
		Values(comment.Name, comment.NewsID, comment.Comment).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
