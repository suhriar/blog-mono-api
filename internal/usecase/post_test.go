package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suhriar/blog-mono-api/internal/repository/mysql/mocks"
	"github.com/suhriar/blog-mono-api/model"
)

func TestCreatePost(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	req := model.CreatePostRequest{
		PostTitle:    "Test Post",
		PostContent:  "This is a test post content.",
		PostHashtags: []string{"golang", "testing"},
	}

	t.Run("Success CreatePost", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("CreatePost", ctx, mock.AnythingOfType("model.Post")).Return(int64(1), nil)

		err := usecase.CreatePost(ctx, userID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail CreatePost - Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("CreatePost", ctx, mock.AnythingOfType("model.Post")).Return(int64(0), assert.AnError)

		err := usecase.CreatePost(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPostByID(t *testing.T) {
	ctx := context.Background()
	postID := int64(1)

	mockPostDetail := model.PostDetail{
		ID:           postID,
		UserID:       1,
		Username:     "testuser",
		PostTitle:    "Test Post",
		PostContent:  "This is a test post content.",
		PostHashtags: []string{"golang", "testing"},
		IsLiked:      true,
	}

	mockComments := []model.CommentResponse{
		{ID: 1, UserID: 2, Username: "commenter", CommentContent: "Nice post!"},
	}

	t.Run("Success GetPostByID", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetPostByID", ctx, postID).Return(mockPostDetail, nil)
		mockRepo.On("CountLikeByPostID", ctx, postID).Return(10, nil)
		mockRepo.On("GetCommentsByPostID", ctx, postID).Return(mockComments, nil)

		post, err := usecase.GetPostByID(ctx, postID)

		assert.NoError(t, err)
		assert.Equal(t, postID, post.PostDetail.ID)
		assert.Equal(t, 10, post.LikeCount)
		assert.Equal(t, mockComments, post.Comments)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail GetPostByID - Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetPostByID", ctx, postID).Return(model.PostDetail{}, assert.AnError)

		post, err := usecase.GetPostByID(ctx, postID)

		assert.Error(t, err)
		assert.Empty(t, post)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail GetPostByID - Count Like Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetPostByID", ctx, postID).Return(mockPostDetail, nil)
		mockRepo.On("CountLikeByPostID", ctx, postID).Return(0, assert.AnError)

		post, err := usecase.GetPostByID(ctx, postID)

		assert.Error(t, err)
		assert.Empty(t, post)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail GetPostByID - Get Comments Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetPostByID", ctx, postID).Return(mockPostDetail, nil)
		mockRepo.On("CountLikeByPostID", ctx, postID).Return(10, nil)
		mockRepo.On("GetCommentsByPostID", ctx, postID).Return([]model.CommentResponse{}, assert.AnError)

		post, err := usecase.GetPostByID(ctx, postID)

		assert.Error(t, err)
		assert.Empty(t, post)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllPost(t *testing.T) {
	ctx := context.Background()
	pageSize := 10
	pageIndex := 1
	limit := pageSize
	offset := pageSize * (pageIndex - 1)

	t.Run("Success GetAllPost", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		expectedPosts := model.GetAllPostResponse{
			Data: []model.PostDetail{
				{ID: 1, UserID: 1, Username: "user1", PostTitle: "Title 1", PostContent: "Content 1"},
				{ID: 2, UserID: 2, Username: "user2", PostTitle: "Title 2", PostContent: "Content 2"},
			},
			Pagination: model.Pagination{
				Limit:  2,
				Offset: 1,
			},
		}

		mockRepo.On("GetAllPost", ctx, limit, offset).Return(expectedPosts, nil)

		posts, err := usecase.GetAllPost(ctx, pageSize, pageIndex)

		assert.NoError(t, err)
		assert.Equal(t, expectedPosts, posts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail GetAllPost - Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetAllPost", ctx, limit, offset).Return(model.GetAllPostResponse{}, assert.AnError)

		posts, err := usecase.GetAllPost(ctx, pageSize, pageIndex)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Empty(t, posts)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateComment(t *testing.T) {
	ctx := context.Background()
	postID := int64(1)
	userID := int64(2)
	req := model.CreateCommentRequest{
		CommentContent: "This is a test comment",
	}

	t.Run("Success CreateComment", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("CreateComment", ctx, mock.AnythingOfType("model.Comment")).Return(int64(1), nil)

		err := usecase.CreateComment(ctx, postID, userID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail CreateComment - Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("CreateComment", ctx, mock.AnythingOfType("model.Comment")).Return(int64(0), assert.AnError)

		err := usecase.CreateComment(ctx, postID, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpsertUserActivity(t *testing.T) {
	ctx := context.Background()
	postID := int64(1)
	userID := int64(1)
	request := model.UserActivityRequest{IsLiked: true}

	t.Run("Success CreateUserActivity", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(model.UserActivity{ID: 0}, nil)
		mockRepo.On("CreateUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(int64(1), nil)

		err := usecase.UpsertUserActivity(ctx, postID, userID, request)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail CreateUserActivity - Never Liked", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}
		request.IsLiked = false

		mockRepo.On("GetUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(model.UserActivity{ID: 0}, nil)

		err := usecase.UpsertUserActivity(ctx, postID, userID, request)

		assert.Error(t, err)
		assert.Equal(t, "never liked this post", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success UpdateUserActivity", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(model.UserActivity{ID: 1}, nil)
		mockRepo.On("UpdateUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(nil)

		err := usecase.UpsertUserActivity(ctx, postID, userID, request)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail UpdateUserActivity - Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.MockPostRepository)
		usecase := &postUsecase{postRepository: mockRepo}

		mockRepo.On("GetUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(model.UserActivity{ID: 1}, nil)
		mockRepo.On("UpdateUserActivity", ctx, mock.AnythingOfType("model.UserActivity")).Return(assert.AnError)

		err := usecase.UpsertUserActivity(ctx, postID, userID, request)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}
