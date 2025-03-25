package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/suhriar/blog-mono-api/model"
)

// Mock user repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(ctx context.Context, email, username string, id int64) (model.User, error) {
	args := m.Called(ctx, email, username, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user model.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetRefreshToken(ctx context.Context, userID int64, now time.Time) (model.RefreshToken, error) {
	args := m.Called(ctx, userID, now)
	return args.Get(0).(model.RefreshToken), args.Error(1)
}

func (m *MockUserRepository) InsertRefreshToken(ctx context.Context, token model.RefreshToken) (int64, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(int64), args.Error(1)
}

// Mock post repository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post model.Post) (int64, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) GetPostByID(ctx context.Context, postID int64) (model.PostDetail, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(model.PostDetail), args.Error(1)
}

func (m *MockPostRepository) CountLikeByPostID(ctx context.Context, postID int64) (int, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockPostRepository) GetCommentsByPostID(ctx context.Context, postID int64) ([]model.CommentResponse, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]model.CommentResponse), args.Error(1)
}

func (m *MockPostRepository) GetAllPost(ctx context.Context, limit, offset int) (model.GetAllPostResponse, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).(model.GetAllPostResponse), args.Error(1)
}

func (m *MockPostRepository) CreateComment(ctx context.Context, comment model.Comment) (int64, error) {
	args := m.Called(ctx, comment)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) GetUserActivity(ctx context.Context, activity model.UserActivity) (model.UserActivity, error) {
	args := m.Called(ctx, activity)
	return args.Get(0).(model.UserActivity), args.Error(1)
}

func (m *MockPostRepository) CreateUserActivity(ctx context.Context, activity model.UserActivity) (int64, error) {
	args := m.Called(ctx, activity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) UpdateUserActivity(ctx context.Context, activity model.UserActivity) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}
