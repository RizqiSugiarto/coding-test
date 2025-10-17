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
	commentDummyID     = "550e8400-e29b-41d4-a716-446655440000"
	sqlInsertComment   = `INSERT INTO comments \(name, news_id, comment\) VALUES \(\$1,\$2,\$3\)`
	commentNewsDummyID = "550e8400-e29b-41d4-a716-44665544125"
	sqlCheckNewsExists = `SELECT EXISTS\(SELECT 1 FROM news WHERE id = \$1\)`
)

func setupCommentMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CommentRepo) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	pg := &postgres.Postgres{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewPostgresCommentRepo(pg)

	return db, mock, repo
}

func TestCommentRepo_Create(t *testing.T) {
	t.Run("success - create comment", func(t *testing.T) {
		db, mock, repo := setupCommentMockDB(t)
		defer db.Close()

		comment := &entity.Comment{
			Name:    "testing comment",
			NewsID:  commentNewsDummyID,
			Comment: "testing comment",
		}

		// Mock news existence check
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(sqlCheckNewsExists).
			WithArgs(comment.NewsID).
			WillReturnRows(rows)

		// Mock insert
		mock.ExpectExec(sqlInsertComment).
			WithArgs(comment.Name, comment.NewsID, comment.Comment).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(context.Background(), comment)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - news not found", func(t *testing.T) {
		db, mock, repo := setupCommentMockDB(t)
		defer db.Close()

		comment := &entity.Comment{
			Name:    "testing comment",
			NewsID:  "non-existent-news-id",
			Comment: "testing comment",
		}

		// Mock news existence check - news does not exist
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery(sqlCheckNewsExists).
			WithArgs(comment.NewsID).
			WillReturnRows(rows)

		err := repo.Create(context.Background(), comment)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - failed to check news existence", func(t *testing.T) {
		db, mock, repo := setupCommentMockDB(t)
		defer db.Close()

		comment := &entity.Comment{
			Name:    "testing comment",
			NewsID:  commentNewsDummyID,
			Comment: "testing comment",
		}

		// Mock news existence check failure
		mock.ExpectQuery(sqlCheckNewsExists).
			WithArgs(comment.NewsID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Create(context.Background(), comment)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - failed insert", func(t *testing.T) {
		db, mock, repo := setupCommentMockDB(t)
		defer db.Close()

		comment := &entity.Comment{
			Name:    "tester",
			NewsID:  commentNewsDummyID,
			Comment: "fail case",
		}

		// Mock news existence check
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(sqlCheckNewsExists).
			WithArgs(comment.NewsID).
			WillReturnRows(rows)

		// Mock insert failure
		mock.ExpectExec(sqlInsertComment).
			WithArgs(comment.Name, comment.NewsID, comment.Comment).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Create(context.Background(), comment)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
