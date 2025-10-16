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
	sqlInsertCategory      = `INSERT INTO categories \(name\) VALUES \(\$1\) RETURNING id, name, created_at, updated_at`
	sqlSelectCategory      = `SELECT id, name, created_at, updated_at FROM categories WHERE id = \$1`
	sqlSelectAllCategories = `SELECT id, name, created_at, updated_at FROM categories ORDER BY id ASC`
	sqlUpdateCategory      = `UPDATE categories SET name = \$1, updated_at = \$2 WHERE id = \$3`
	sqlDeleteCategory      = `DELETE FROM categories WHERE id = \$1`
)

const (
	dummyID       = "550e8400-e29b-41d4-a716-446655440000"
	nonExistentID = "550e8400-e29b-41d4-a716-999999999999"
)

func setupCategoryMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CategoryRepo) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	pg := &postgres.Postgres{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewPostgresCategoryRepo(pg)

	return db, mock, repo
}

func TestCategoryRepo_Create(t *testing.T) {
	t.Run("success - create category", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		category := &entity.Category{
			Name: "Technology",
		}

		now := time.Now()
		expectedCategory := &entity.Category{
			ID:        dummyID,
			Name:      "Technology",
			CreatedAt: now,
			UpdatedAt: now,
		}

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(expectedCategory.ID, expectedCategory.Name, expectedCategory.CreatedAt, expectedCategory.UpdatedAt)

		mock.ExpectQuery(sqlInsertCategory).
			WithArgs(category.Name).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), category)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCategory.ID, result.ID)
		assert.Equal(t, expectedCategory.Name, result.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database insert fails", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		category := &entity.Category{
			Name: "Technology",
		}

		mock.ExpectQuery(sqlInsertCategory).
			WithArgs(category.Name).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.Create(context.Background(), category)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - create category with special characters", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		category := &entity.Category{
			Name: "Tech & Innovation",
		}

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(dummyID, category.Name, now, now)

		mock.ExpectQuery(sqlInsertCategory).
			WithArgs(category.Name).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), category)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Tech & Innovation", result.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryRepo_GetByID(t *testing.T) {
	t.Run("success - get category by id", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		now := time.Now()
		expectedCategory := &entity.Category{
			ID:        dummyID,
			Name:      "Technology",
			CreatedAt: now,
			UpdatedAt: now,
		}

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(expectedCategory.ID, expectedCategory.Name, expectedCategory.CreatedAt, expectedCategory.UpdatedAt)

		mock.ExpectQuery(sqlSelectCategory).
			WithArgs(expectedCategory.ID).
			WillReturnRows(rows)

		result, err := repo.GetByID(context.Background(), expectedCategory.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCategory.ID, result.ID)
		assert.Equal(t, expectedCategory.Name, result.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - category not found", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		categoryID := nonExistentID

		mock.ExpectQuery(sqlSelectCategory).
			WithArgs(categoryID).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(context.Background(), categoryID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		categoryID := dummyID

		mock.ExpectQuery(sqlSelectCategory).
			WithArgs(categoryID).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.GetByID(context.Background(), categoryID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryRepo_GetAll(t *testing.T) {
	t.Run("success - get all categories", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440001", "Technology", now, now).
			AddRow("550e8400-e29b-41d4-a716-446655440002", "Sports", now, now).
			AddRow("550e8400-e29b-41d4-a716-446655440003", "Entertainment", now, now)

		mock.ExpectQuery(sqlSelectAllCategories).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, "Technology", result[0].Name)
		assert.Equal(t, "Sports", result[1].Name)
		assert.Equal(t, "Entertainment", result[2].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - get all categories empty result", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})

		mock.ExpectQuery(sqlSelectAllCategories).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		mock.ExpectQuery(sqlSelectAllCategories).
			WillReturnError(apperror.ErrDatabaseConnection)

		result, err := repo.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryRepo_Update(t *testing.T) {
	t.Run("success - update category", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		category := &entity.Category{
			ID:   dummyID,
			Name: "Updated Technology",
		}

		mock.ExpectExec(sqlUpdateCategory).
			WithArgs(category.Name, "NOW()", category.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(context.Background(), category)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - category not found", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		category := &entity.Category{
			ID:   nonExistentID,
			Name: "Non-existent Category",
		}

		mock.ExpectExec(sqlUpdateCategory).
			WithArgs(category.Name, "NOW()", category.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(context.Background(), category)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database update fails", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		category := &entity.Category{
			ID:   dummyID,
			Name: "Updated Technology",
		}

		mock.ExpectExec(sqlUpdateCategory).
			WithArgs(category.Name, "NOW()", category.ID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Update(context.Background(), category)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryRepo_Delete(t *testing.T) {
	t.Run("success - delete category", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		categoryID := dummyID

		mock.ExpectExec(sqlDeleteCategory).
			WithArgs(categoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(context.Background(), categoryID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - category not found", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		categoryID := nonExistentID

		mock.ExpectExec(sqlDeleteCategory).
			WithArgs(categoryID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(context.Background(), categoryID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database delete fails", func(t *testing.T) {
		db, mock, repo := setupCategoryMockDB(t)
		defer db.Close()

		categoryID := dummyID

		mock.ExpectExec(sqlDeleteCategory).
			WithArgs(categoryID).
			WillReturnError(apperror.ErrDatabaseConnection)

		err := repo.Delete(context.Background(), categoryID)

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrDatabaseConnection, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewPostgresCategoryRepo(t *testing.T) {
	t.Run("success - create new category repository", func(t *testing.T) {
		db, _, _ := setupCategoryMockDB(t)
		defer db.Close()

		pg := &postgres.Postgres{
			DB:      db,
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		}

		repo := NewPostgresCategoryRepo(pg)

		assert.NotNil(t, repo)
		assert.NotNil(t, repo.Postgres)
		assert.NotNil(t, repo.DB)
		assert.NotNil(t, repo.Builder)
	})
}
