package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepo struct {
	mock.Mock
}

const (
	testCommentID      = "550e8400-e29b-41d4-a716-446655440000"
	testCommentNewsID  = "550e8400-e29b-41d4-a716-446655440000"
	testCommentName    = "John Doe"
	testCommentContent = "This is a great article!"
)

var errDatabaseError = errors.New("database error")

func (m *MockCommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)

	return args.Error(0)
}

func TestCommentUseCase_Create(t *testing.T) {
	t.Run("success - create comment", func(t *testing.T) {
		mockRepo := new(MockCommentRepo)
		mockUseCase := NewCommentUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCommentRequestDTO{
			Name:    testCommentName,
			Comment: testCommentContent,
			NewsID:  testCommentNewsID,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(comment *entity.Comment) bool {
			return comment.Name == testCommentName &&
				comment.Comment == testCommentContent &&
				comment.NewsID == testCommentNewsID
		})).Return(nil)

		err := mockUseCase.Create(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository create fails", func(t *testing.T) {
		mockRepo := new(MockCommentRepo)
		mockUseCase := NewCommentUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCommentRequestDTO{
			Name:    testCommentName,
			Comment: testCommentContent,
			NewsID:  testCommentNewsID,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(comment *entity.Comment) bool {
			return comment.Name == testCommentName &&
				comment.Comment == testCommentContent &&
				comment.NewsID == testCommentNewsID
		})).Return(errDatabaseError)

		err := mockUseCase.Create(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, errDatabaseError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - create comment with empty name", func(t *testing.T) {
		mockRepo := new(MockCommentRepo)
		mockUseCase := NewCommentUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCommentRequestDTO{
			Name:    "",
			Comment: "Anonymous comment",
			NewsID:  testCommentNewsID,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(comment *entity.Comment) bool {
			return comment.Name == "" &&
				comment.Comment == "Anonymous comment" &&
				comment.NewsID == testCommentNewsID
		})).Return(nil)

		err := mockUseCase.Create(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - create comment with special characters", func(t *testing.T) {
		mockRepo := new(MockCommentRepo)
		mockUseCase := NewCommentUseCase(mockRepo)

		ctx := context.Background()
		req := &dto.CreateCommentRequestDTO{
			Name:    "Test User <script>",
			Comment: "Comment with special chars: !@#$%^&*()",
			NewsID:  testCommentNewsID,
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(comment *entity.Comment) bool {
			return comment.Name == "Test User <script>" &&
				comment.Comment == "Comment with special chars: !@#$%^&*()" &&
				comment.NewsID == testCommentNewsID
		})).Return(nil)

		err := mockUseCase.Create(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
