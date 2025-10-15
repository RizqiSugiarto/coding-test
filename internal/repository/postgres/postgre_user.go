package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewPostgresUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) Create(ctx context.Context, user entity.User) error {
	query, args, err := u.Builder.Insert("users").
		Columns("username, password").
		Values(user.Username, user.Password).
		ToSql()
	if err != nil {
		return err
	}

	_, err = u.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query, args, err := u.Builder.
		Select("id", "username", "password").
		From("users").
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := u.DB.QueryRowContext(ctx, query, args...)

	var user entity.User

	err = row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}
