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

type CustomPageRepo struct {
	*postgres.Postgres
}

func NewPostgresCustomPageRepo(pg *postgres.Postgres) *CustomPageRepo {
	return &CustomPageRepo{pg}
}

func (r *CustomPageRepo) Create(ctx context.Context, page *entity.CustomPage) (*entity.CustomPage, error) {
	query := r.Builder.
		Insert("custom_pages").
		Columns("custom_url", "content", "author_id").
		Values(page.CustomURL, page.Content, page.AuthorID).
		Suffix("RETURNING id, custom_url, content, author_id, created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var result entity.CustomPage

	err = r.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&result.ID,
		&result.CustomURL,
		&result.Content,
		&result.AuthorID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *CustomPageRepo) GetByID(ctx context.Context, id string) (*entity.CustomPage, error) {
	query := r.Builder.
		Select("id", "custom_url", "content", "author_id", "created_at", "updated_at").
		From("custom_pages").
		Where(squirrel.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var page entity.CustomPage

	err = r.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&page.ID,
		&page.CustomURL,
		&page.Content,
		&page.AuthorID,
		&page.CreatedAt,
		&page.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}

		return nil, err
	}

	return &page, nil
}

func (r *CustomPageRepo) GetAll(ctx context.Context) ([]entity.CustomPage, error) {
	query := r.Builder.
		Select("id", "custom_url", "content", "author_id", "created_at", "updated_at").
		From("custom_pages").
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

	var pages []entity.CustomPage

	for rows.Next() {
		var page entity.CustomPage

		err := rows.Scan(
			&page.ID,
			&page.CustomURL,
			&page.Content,
			&page.AuthorID,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if pages == nil {
		pages = []entity.CustomPage{}
	}

	return pages, nil
}

func (r *CustomPageRepo) Update(ctx context.Context, page *entity.CustomPage) error {
	query := r.Builder.
		Update("custom_pages").
		Set("custom_url", page.CustomURL).
		Set("content", page.Content).
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": page.ID})

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

func (r *CustomPageRepo) Delete(ctx context.Context, id string) error {
	query := r.Builder.
		Delete("custom_pages").
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
