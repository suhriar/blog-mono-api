package mysql

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/suhriar/blog-mono-api/model"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}

	ctx := context.Background()
	post := model.Post{
		UserID:       1,
		PostTitle:    "Test Title",
		PostContent:  "Test Content",
		PostHashtags: "test,go,sqlmock",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    "test_user",
		UpdatedBy:    "test_user",
	}

	mock.ExpectExec(`INSERT INTO posts`).
		WithArgs(post.UserID, post.PostTitle, post.PostContent, post.PostHashtags, post.CreatedAt, post.UpdatedAt, post.CreatedBy, post.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	lastInsertID, err := repo.CreatePost(ctx, post)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}

	ctx := context.Background()
	limit, offset := 10, 0

	expectedPosts := []model.PostDetail{
		{ID: 1, UserID: 2, Username: "user1", PostTitle: "Title 1", PostContent: "Content 1", PostHashtags: []string{"tag1", "tag2"}},
		{ID: 2, UserID: 3, Username: "user2", PostTitle: "Title 2", PostContent: "Content 2", PostHashtags: []string{"tag3", "tag4"}},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "username", "post_title", "post_content", "post_hashtags"}).
		AddRow(expectedPosts[0].ID, expectedPosts[0].UserID, expectedPosts[0].Username, expectedPosts[0].PostTitle, expectedPosts[0].PostContent, strings.Join(expectedPosts[0].PostHashtags, ",")).
		AddRow(expectedPosts[1].ID, expectedPosts[1].UserID, expectedPosts[1].Username, expectedPosts[1].PostTitle, expectedPosts[1].PostContent, strings.Join(expectedPosts[1].PostHashtags, ","))

	mock.ExpectQuery(`SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags FROM posts p JOIN users u ON p.user_id = u.id`).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	resp, err := repo.GetAllPost(ctx, limit, offset)
	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, resp.Data)
	assert.Equal(t, limit, resp.Pagination.Limit)
	assert.Equal(t, offset, resp.Pagination.Offset)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &postRepository{db: db}

	ctx := context.Background()
	postID := int64(1)

	expectedPost := model.PostDetail{
		ID:           1,
		UserID:       2,
		Username:     "user1",
		PostTitle:    "Title 1",
		PostContent:  "Content 1",
		PostHashtags: []string{"tag1", "tag2"},
		IsLiked:      true,
	}

	row := sqlmock.NewRows([]string{"id", "user_id", "username", "post_title", "post_content", "post_hashtags", "is_liked"}).
		AddRow(expectedPost.ID, expectedPost.UserID, expectedPost.Username, expectedPost.PostTitle, expectedPost.PostContent, strings.Join(expectedPost.PostHashtags, ","), expectedPost.IsLiked)

	mock.ExpectQuery(`SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, uv.is_liked FROM posts p JOIN users u ON p.user_id = u.id`).
		WithArgs(postID).
		WillReturnRows(row)

	resp, err := repo.GetPostByID(ctx, postID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, resp)
	assert.NoError(t, mock.ExpectationsWereMet())
}
