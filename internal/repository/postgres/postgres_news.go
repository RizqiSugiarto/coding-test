package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/postgres"
)

// NewsRepo implements repository.NewsRepo interface.
type NewsRepo struct {
	*postgres.Postgres
}

// NewPostgresNewsRepo creates a new PostgreSQL news repository.
func NewPostgresNewsRepo(pg *postgres.Postgres) *NewsRepo {
	return &NewsRepo{pg}
}

func (r *NewsRepo) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	query := r.Builder.
		Insert("news").
		Columns("category_id", "author_id", "title", "content").
		Values(news.CategoryID, news.AuthorID, news.Title, news.Content).
		Suffix("RETURNING id, category_id, author_id, title, content, created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var result entity.News

	err = r.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&result.ID,
		&result.CategoryID,
		&result.AuthorID,
		&result.Title,
		&result.Content,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *NewsRepo) GetByID(ctx context.Context, id string) (*entity.News, error) {
	query := r.Builder.
		Select("id", "category_id", "author_id", "title", "content", "created_at", "updated_at").
		From("news").
		Where(squirrel.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var news entity.News

	err = r.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&news.ID,
		&news.CategoryID,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.CreatedAt,
		&news.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}

		return nil, err
	}

	return &news, nil
}

func (r *NewsRepo) GetAll(ctx context.Context) ([]entity.News, error) {
	query := r.Builder.
		Select("id", "category_id", "author_id", "title", "content", "created_at", "updated_at").
		From("news").
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []entity.News

	for rows.Next() {
		var news entity.News

		err := rows.Scan(
			&news.ID,
			&news.CategoryID,
			&news.AuthorID,
			&news.Title,
			&news.Content,
			&news.CreatedAt,
			&news.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if newsList == nil {
		newsList = []entity.News{}
	}

	return newsList, nil
}

func (r *NewsRepo) Update(ctx context.Context, news *entity.News) error {
	query := r.Builder.
		Update("news").
		Set("category_id", news.CategoryID).
		Set("title", news.Title).
		Set("content", news.Content).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": news.ID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *NewsRepo) Delete(ctx context.Context, id string) error {
	query := r.Builder.
		Delete("news").
		Where(squirrel.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}
