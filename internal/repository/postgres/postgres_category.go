package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/postgres"
)

type CategoryRepo struct {
	*postgres.Postgres
}

func NewPostgresCategoryRepo(pg *postgres.Postgres) *CategoryRepo {
	return &CategoryRepo{pg}
}

func (c *CategoryRepo) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	query, args, err := c.Builder.Insert("categories").
		Columns("name").
		Values(category.Name).
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var result entity.Category

	err = c.DB.QueryRowContext(ctx, query, args...).Scan(
		&result.ID,
		&result.Name,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *CategoryRepo) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	query, args, err := c.Builder.
		Select("id", "name", "created_at", "updated_at").
		From("categories").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := c.DB.QueryRowContext(ctx, query, args...)

	var category entity.Category

	err = row.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}

		return nil, err
	}

	return &category, nil
}

func (c *CategoryRepo) GetAll(ctx context.Context) ([]entity.Category, error) {
	query, args, err := c.Builder.
		Select("id", "name", "created_at", "updated_at").
		From("categories").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]entity.Category, 0)

	for rows.Next() {
		var category entity.Category

		err = rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoryRepo) Update(ctx context.Context, category *entity.Category) error {
	query, args, err := c.Builder.Update("categories").
		Set("name", category.Name).
		Set("updated_at", "NOW()").
		Where("id = ?", category.ID).
		ToSql()
	if err != nil {
		return err
	}

	result, err := c.DB.ExecContext(ctx, query, args...)
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

func (c *CategoryRepo) Delete(ctx context.Context, id string) error {
	query, args, err := c.Builder.Delete("categories").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return err
	}

	result, err := c.DB.ExecContext(ctx, query, args...)
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
