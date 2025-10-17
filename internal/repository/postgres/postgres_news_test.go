package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sqlInsertNews     = `INSERT INTO news \(category_id,author_id,title,content\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING id, category_id, author_id, title, content, created_at, updated_at`
	sqlSelectNews     = `SELECT id, category_id, author_id, title, content, created_at, updated_at FROM news WHERE id = \$1`
	sqlSelectAllNews  = `SELECT id, category_id, author_id, title, content, created_at, updated_at FROM news ORDER BY created_at DESC`
	sqlUpdateNews     = `UPDATE news SET category_id = \$1, title = \$2, content = \$3, updated_at = NOW\(\) WHERE id = \$4`
	sqlDeleteNews     = `DELETE FROM news WHERE id = \$1`
	testNewsID        = "550e8400-e29b-41d4-a716-446655440000"
	testCategoryID    = "550e8400-e29b-41d4-a716-446655440001"
	testAuthorID      = "550e8400-e29b-41d4-a716-446655440002"
	nonExistentNewsID = "550e8400-e29b-41d4-a716-999999999999"
)

func setupNewsMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *NewsRepo) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	pg := &postgres.Postgres{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewPostgresNewsRepo(pg)

	return db, mock, repo
}

func TestNewsRepo_Create(t *testing.T) {
	t.Run("success - create news", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		news := &entity.News{
			CategoryID: testCategoryID,
			AuthorID:   testAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
		}

		now := time.Now()
		expectedNews := &entity.News{
			ID:         testNewsID,
			CategoryID: testCategoryID,
			AuthorID:   testAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{"id", "category_id", "author_id", "title", "content", "created_at", "updated_at"}).
			AddRow(expectedNews.ID, expectedNews.CategoryID, expectedNews.AuthorID, expectedNews.Title, expectedNews.Content, expectedNews.CreatedAt, expectedNews.UpdatedAt)

		mock.ExpectQuery(sqlInsertNews).
			WithArgs(news.CategoryID, news.AuthorID, news.Title, news.Content).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), news)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedNews.ID, result.ID)
		assert.Equal(t, expectedNews.Title, result.Title)
		assert.Equal(t, expectedNews.Content, result.Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database insert fails", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		news := &entity.News{
			CategoryID: testCategoryID,
			AuthorID:   testAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
		}

		mock.ExpectQuery(sqlInsertNews).
			WithArgs(news.CategoryID, news.AuthorID, news.Title, news.Content).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.Create(context.Background(), news)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - create news with long content", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		longContent := "This is a very long content that spans multiple lines and contains detailed information about the news article."
		news := &entity.News{
			CategoryID: testCategoryID,
			AuthorID:   testAuthorID,
			Title:      "Breaking News",
			Content:    longContent,
		}

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "category_id", "author_id", "title", "content", "created_at", "updated_at"}).
			AddRow(testNewsID, testCategoryID, testAuthorID, news.Title, longContent, now, now)

		mock.ExpectQuery(sqlInsertNews).
			WithArgs(news.CategoryID, news.AuthorID, news.Title, news.Content).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), news)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, longContent, result.Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewsRepo_GetByID(t *testing.T) {
	t.Run("success - get news by id", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		now := time.Now()
		expectedNews := &entity.News{
			ID:         testNewsID,
			CategoryID: testCategoryID,
			AuthorID:   testAuthorID,
			Title:      "Breaking News",
			Content:    "This is the news content",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{"id", "category_id", "author_id", "title", "content", "created_at", "updated_at"}).
			AddRow(expectedNews.ID, expectedNews.CategoryID, expectedNews.AuthorID, expectedNews.Title, expectedNews.Content, expectedNews.CreatedAt, expectedNews.UpdatedAt)

		mock.ExpectQuery(sqlSelectNews).
			WithArgs(expectedNews.ID).
			WillReturnRows(rows)

		result, err := repo.GetByID(context.Background(), expectedNews.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedNews.ID, result.ID)
		assert.Equal(t, expectedNews.Title, result.Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - news not found", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectNews).
			WithArgs(nonExistentNewsID).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(context.Background(), nonExistentNewsID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectNews).
			WithArgs(testNewsID).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.GetByID(context.Background(), testNewsID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewsRepo_GetAll(t *testing.T) {
	t.Run("success - get all news", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "category_id", "author_id", "title", "content", "created_at", "updated_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440001", testCategoryID, testAuthorID, "News 1", "Content 1", now, now).
			AddRow("550e8400-e29b-41d4-a716-446655440002", testCategoryID, testAuthorID, "News 2", "Content 2", now, now).
			AddRow("550e8400-e29b-41d4-a716-446655440003", testCategoryID, testAuthorID, "News 3", "Content 3", now, now)

		mock.ExpectQuery(sqlSelectAllNews).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, "News 1", result[0].Title)
		assert.Equal(t, "News 2", result[1].Title)
		assert.Equal(t, "News 3", result[2].Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - get all news empty result", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "category_id", "author_id", "title", "content", "created_at", "updated_at"})

		mock.ExpectQuery(sqlSelectAllNews).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectAllNews).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewsRepo_Update(t *testing.T) {
	t.Run("success - update news", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		news := &entity.News{
			ID:         testNewsID,
			CategoryID: testCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}

		mock.ExpectExec(sqlUpdateNews).
			WithArgs(news.CategoryID, news.Title, news.Content, news.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(context.Background(), news)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - news not found", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		news := &entity.News{
			ID:         nonExistentNewsID,
			CategoryID: testCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}

		mock.ExpectExec(sqlUpdateNews).
			WithArgs(news.CategoryID, news.Title, news.Content, news.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(context.Background(), news)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database update fails", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		news := &entity.News{
			ID:         testNewsID,
			CategoryID: testCategoryID,
			Title:      "Updated News",
			Content:    "Updated content",
		}

		mock.ExpectExec(sqlUpdateNews).
			WithArgs(news.CategoryID, news.Title, news.Content, news.ID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Update(context.Background(), news)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewsRepo_Delete(t *testing.T) {
	t.Run("success - delete news", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		mock.ExpectExec(sqlDeleteNews).
			WithArgs(testNewsID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(context.Background(), testNewsID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - news not found", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		mock.ExpectExec(sqlDeleteNews).
			WithArgs(nonExistentNewsID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(context.Background(), nonExistentNewsID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database delete fails", func(t *testing.T) {
		db, mock, repo := setupNewsMockDB(t)
		defer db.Close()

		mock.ExpectExec(sqlDeleteNews).
			WithArgs(testNewsID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Delete(context.Background(), testNewsID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
