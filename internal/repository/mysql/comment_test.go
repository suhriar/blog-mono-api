package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/suhriar/blog-mono-api/model"
)

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}

	ctx := context.Background()
	comment := model.Comment{
		PostID:         1,
		UserID:         2,
		CommentContent: "This is a test comment",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      "test_user",
		UpdatedBy:      "test_user",
	}

	mock.ExpectExec(`INSERT INTO comments`).
		WithArgs(comment.PostID, comment.UserID, comment.CommentContent, comment.CreatedAt, comment.UpdatedAt, comment.CreatedBy, comment.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	lastInsertID, err := repo.CreateComment(ctx, comment)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCommentsByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}

	ctx := context.Background()
	postID := int64(1)

	expectedComments := []model.CommentResponse{
		{ID: 1, UserID: 2, CommentContent: "This is a test comment", Username: "test_user"},
		{ID: 2, UserID: 3, CommentContent: "Another test comment", Username: "another_user"},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "comment_content", "username"}).
		AddRow(expectedComments[0].ID, expectedComments[0].UserID, expectedComments[0].CommentContent, expectedComments[0].Username).
		AddRow(expectedComments[1].ID, expectedComments[1].UserID, expectedComments[1].CommentContent, expectedComments[1].Username)

	mock.ExpectQuery(`SELECT c.id, c.user_id, c.comment_content, u.username FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = \?`).
		WithArgs(postID).
		WillReturnRows(rows)

	comments, err := repo.GetCommentsByPostID(ctx, postID)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	assert.NoError(t, mock.ExpectationsWereMet())
}
