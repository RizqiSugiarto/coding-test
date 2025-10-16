package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sqlInsertUser           = `INSERT INTO users \(username, password\) VALUES \(\$1,\$2\)`
	sqlSelectUserByUsername = `SELECT id, username, password FROM users WHERE username = \$1`
)

// setupMockDB creates a mock database and returns the mock controller.
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *UserRepo) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	pg := &postgres.Postgres{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewPostgresUserRepo(pg)

	return db, mock, repo
}

func TestUserRepo_Create(t *testing.T) {
	t.Run("success - create user", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		user := entity.User{
			Username: "testuser",
			Password: "hashedpassword123",
		}

		expectedSQL := sqlInsertUser
		mock.ExpectExec(expectedSQL).
			WithArgs(user.Username, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(context.Background(), user)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database exec fails", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		user := entity.User{
			Username: "testuser",
			Password: "hashedpassword123",
		}

		expectedSQL := sqlInsertUser
		mock.ExpectExec(expectedSQL).
			WithArgs(user.Username, user.Password).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Create(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - duplicate username constraint violation", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		user := entity.User{
			Username: "existinguser",
			Password: "hashedpassword123",
		}

		expectedSQL := sqlInsertUser
		mock.ExpectExec(expectedSQL).
			WithArgs(user.Username, user.Password).
			WillReturnError(apperror.ErrDuplicateKey)

		err := repo.Create(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - create user with empty password", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		user := entity.User{
			Username: "testuser",
			Password: "",
		}

		expectedSQL := sqlInsertUser
		mock.ExpectExec(expectedSQL).
			WithArgs(user.Username, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(context.Background(), user)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepo_GetByUsername(t *testing.T) {
	t.Run("success - get user by username", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		expectedUser := &entity.User{
			ID:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "testuser",
			Password: "hashedpassword123",
		}

		expectedSQL := sqlSelectUserByUsername
		rows := sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password)

		mock.ExpectQuery(expectedSQL).
			WithArgs(expectedUser.Username).
			WillReturnRows(rows)

		user, err := repo.GetByUsername(context.Background(), expectedUser.Username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Password, user.Password)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - user not found", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		username := "nonexistentuser"

		expectedSQL := sqlSelectUserByUsername
		mock.ExpectQuery(expectedSQL).
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByUsername(context.Background(), username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		username := "testuser"

		expectedSQL := sqlSelectUserByUsername
		mock.ExpectQuery(expectedSQL).
			WithArgs(username).
			WillReturnError(apperror.ErrDatabaseConnection)

		user, err := repo.GetByUsername(context.Background(), username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - get user with special characters in username", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		expectedUser := &entity.User{
			ID:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "test.user+special@domain",
			Password: "hashedpassword123",
		}

		expectedSQL := sqlSelectUserByUsername
		rows := sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password)

		mock.ExpectQuery(expectedSQL).
			WithArgs(expectedUser.Username).
			WillReturnRows(rows)

		user, err := repo.GetByUsername(context.Background(), expectedUser.Username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - empty username", func(t *testing.T) {
		db, mock, repo := setupMockDB(t)
		defer db.Close()

		username := ""

		expectedSQL := sqlSelectUserByUsername
		mock.ExpectQuery(expectedSQL).
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByUsername(context.Background(), username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestNewPostgresUserRepo tests the repository constructor.
func TestNewPostgresUserRepo(t *testing.T) {
	t.Run("success - create new repository", func(t *testing.T) {
		db, _, _ := setupMockDB(t)
		defer db.Close()

		pg := &postgres.Postgres{
			DB:      db,
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		}

		repo := NewPostgresUserRepo(pg)

		assert.NotNil(t, repo)
		assert.NotNil(t, repo.Postgres)
		assert.NotNil(t, repo.DB)
		assert.NotNil(t, repo.Builder)
	})
}
