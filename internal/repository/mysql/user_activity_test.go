package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/suhriar/blog-mono-api/model"
)

func TestGetUserActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}
	ctx := context.Background()
	userActivity := model.UserActivity{PostID: 1, UserID: 2}

	mock.ExpectQuery(`SELECT id, post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by FROM user_activities WHERE post_id = \? AND user_id = \?`).
		WithArgs(userActivity.PostID, userActivity.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "is_liked", "created_at", "updated_at", "created_by", "updated_by"}).
			AddRow(1, 1, 2, true, time.Now(), time.Now(), "test_user", "test_user"))

	resp, err := repo.GetUserActivity(ctx, userActivity)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}
	ctx := context.Background()
	userActivity := model.UserActivity{PostID: 1, UserID: 2, IsLiked: true, CreatedAt: time.Now(), UpdatedAt: time.Now(), CreatedBy: "test_user", UpdatedBy: "test_user"}

	mock.ExpectExec(`INSERT INTO user_activities`).
		WithArgs(userActivity.PostID, userActivity.UserID, userActivity.IsLiked, userActivity.CreatedAt, userActivity.UpdatedAt, userActivity.CreatedBy, userActivity.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	lastInsertID, err := repo.CreateUserActivity(ctx, userActivity)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}
	ctx := context.Background()
	userActivity := model.UserActivity{PostID: 1, UserID: 2, IsLiked: false, UpdatedAt: time.Now(), UpdatedBy: "test_user"}

	mock.ExpectExec(`UPDATE user_activities SET is_liked = \?, updated_at = \?, updated_by = \? WHERE post_id = \? AND user_id = \?`).
		WithArgs(userActivity.IsLiked, userActivity.UpdatedAt, userActivity.UpdatedBy, userActivity.PostID, userActivity.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUserActivity(ctx, userActivity)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCountLikeByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}
	ctx := context.Background()
	postID := int64(1)

	mock.ExpectQuery(`SELECT COUNT\(id\) FROM user_activities WHERE post_id = \? AND is_liked = true`).
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.CountLikeByPostID(ctx, postID)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
