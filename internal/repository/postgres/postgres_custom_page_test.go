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
	sqlInsertPage        = `INSERT INTO custom_pages \(custom_url,content,author_id\) VALUES \(\$1,\$2,\$3\) RETURNING id, custom_url, content, author_id, created_at, updated_at`
	sqlSelectPage        = `SELECT id, custom_url, content, author_id, created_at, updated_at FROM custom_pages WHERE id = \$1`
	sqlSelectAllPages    = `SELECT id, custom_url, content, author_id, created_at, updated_at FROM custom_pages ORDER BY created_at DESC`
	sqlUpdatePage        = `UPDATE custom_pages SET custom_url = \$1, content = \$2, updated_at = CURRENT_TIMESTAMP WHERE id = \$3`
	sqlDeletePage        = `DELETE FROM custom_pages WHERE id = \$1`
	testPageID           = "550e8400-e29b-41d4-a716-446655440000"
	testPageAuthorID     = "550e8400-e29b-41d4-a716-446655440001"
	nonExistentPageID    = "550e8400-e29b-41d4-a716-999999999999"
	testCustomURL        = "/about-us"
	testCustomURLUpdated = "/about-company"
)

func setupPageMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CustomPageRepo) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	pg := &postgres.Postgres{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewPostgresCustomPageRepo(pg)

	return db, mock, repo
}

func TestCustomPageRepo_Create(t *testing.T) {
	t.Run("success - create custom page", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		page := &entity.CustomPage{
			CustomURL: testCustomURL,
			Content:   "This is the about us page content",
			AuthorID:  testPageAuthorID,
		}

		now := time.Now()
		expectedPage := &entity.CustomPage{
			ID:        testPageID,
			CustomURL: testCustomURL,
			Content:   "This is the about us page content",
			AuthorID:  testPageAuthorID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		rows := sqlmock.NewRows([]string{"id", "custom_url", "content", "author_id", "created_at", "updated_at"}).
			AddRow(expectedPage.ID, expectedPage.CustomURL, expectedPage.Content, expectedPage.AuthorID, expectedPage.CreatedAt, expectedPage.UpdatedAt)

		mock.ExpectQuery(sqlInsertPage).
			WithArgs(page.CustomURL, page.Content, page.AuthorID).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), page)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedPage.ID, result.ID)
		assert.Equal(t, expectedPage.CustomURL, result.CustomURL)
		assert.Equal(t, expectedPage.Content, result.Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database insert fails", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		page := &entity.CustomPage{
			CustomURL: testCustomURL,
			Content:   "This is the about us page content",
			AuthorID:  testPageAuthorID,
		}

		mock.ExpectQuery(sqlInsertPage).
			WithArgs(page.CustomURL, page.Content, page.AuthorID).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.Create(context.Background(), page)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - create custom page with long content", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		longContent := "This is a very long content that spans multiple lines and contains detailed information about the custom page. It includes HTML markup, CSS styles, and JavaScript code for a fully functional page."
		page := &entity.CustomPage{
			CustomURL: testCustomURL,
			Content:   longContent,
			AuthorID:  testPageAuthorID,
		}

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "custom_url", "content", "author_id", "created_at", "updated_at"}).
			AddRow(testPageID, testCustomURL, longContent, testPageAuthorID, now, now)

		mock.ExpectQuery(sqlInsertPage).
			WithArgs(page.CustomURL, page.Content, page.AuthorID).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), page)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, longContent, result.Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCustomPageRepo_GetByID(t *testing.T) {
	t.Run("success - get custom page by id", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		now := time.Now()
		expectedPage := &entity.CustomPage{
			ID:        testPageID,
			CustomURL: testCustomURL,
			Content:   "This is the about us page content",
			AuthorID:  testPageAuthorID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		rows := sqlmock.NewRows([]string{"id", "custom_url", "content", "author_id", "created_at", "updated_at"}).
			AddRow(expectedPage.ID, expectedPage.CustomURL, expectedPage.Content, expectedPage.AuthorID, expectedPage.CreatedAt, expectedPage.UpdatedAt)

		mock.ExpectQuery(sqlSelectPage).
			WithArgs(expectedPage.ID).
			WillReturnRows(rows)

		result, err := repo.GetByID(context.Background(), expectedPage.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedPage.ID, result.ID)
		assert.Equal(t, expectedPage.CustomURL, result.CustomURL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - custom page not found", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectPage).
			WithArgs(nonExistentPageID).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(context.Background(), nonExistentPageID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectPage).
			WithArgs(testPageID).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.GetByID(context.Background(), testPageID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCustomPageRepo_GetAll(t *testing.T) {
	t.Run("success - get all custom pages", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "custom_url", "content", "author_id", "created_at", "updated_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440001", "/about-us", "About content", testPageAuthorID, now, now).
			AddRow("550e8400-e29b-41d4-a716-446655440002", "/contact", "Contact content", testPageAuthorID, now, now).
			AddRow("550e8400-e29b-41d4-a716-446655440003", "/privacy-policy", "Privacy content", testPageAuthorID, now, now)

		mock.ExpectQuery(sqlSelectAllPages).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, "/about-us", result[0].CustomURL)
		assert.Equal(t, "/contact", result[1].CustomURL)
		assert.Equal(t, "/privacy-policy", result[2].CustomURL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - get all custom pages empty result", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "custom_url", "content", "author_id", "created_at", "updated_at"})

		mock.ExpectQuery(sqlSelectAllPages).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectAllPages).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCustomPageRepo_Update(t *testing.T) {
	t.Run("success - update custom page", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		page := &entity.CustomPage{
			ID:        testPageID,
			CustomURL: testCustomURLUpdated,
			Content:   "Updated content",
		}

		mock.ExpectExec(sqlUpdatePage).
			WithArgs(page.CustomURL, page.Content, page.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(context.Background(), page)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - custom page not found", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		page := &entity.CustomPage{
			ID:        nonExistentPageID,
			CustomURL: testCustomURLUpdated,
			Content:   "Updated content",
		}

		mock.ExpectExec(sqlUpdatePage).
			WithArgs(page.CustomURL, page.Content, page.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(context.Background(), page)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database update fails", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		page := &entity.CustomPage{
			ID:        testPageID,
			CustomURL: testCustomURLUpdated,
			Content:   "Updated content",
		}

		mock.ExpectExec(sqlUpdatePage).
			WithArgs(page.CustomURL, page.Content, page.ID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Update(context.Background(), page)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCustomPageRepo_Delete(t *testing.T) {
	t.Run("success - delete custom page", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		mock.ExpectExec(sqlDeletePage).
			WithArgs(testPageID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(context.Background(), testPageID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - custom page not found", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		mock.ExpectExec(sqlDeletePage).
			WithArgs(nonExistentPageID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(context.Background(), nonExistentPageID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database delete fails", func(t *testing.T) {
		db, mock, repo := setupPageMockDB(t)
		defer db.Close()

		mock.ExpectExec(sqlDeletePage).
			WithArgs(testPageID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Delete(context.Background(), testPageID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
