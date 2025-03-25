package usecase

import (
	"context"

	repository "github.com/suhriar/blog-mono-api/internal/repository/mysql"
	"github.com/suhriar/blog-mono-api/model"
)

type UserUsecase interface {
	SignUp(ctx context.Context, req model.SignUpRequest) (err error)
	Login(ctx context.Context, req model.LoginRequest) (jwtToken, refreshToken string, err error)
	ValidateRefreshToken(ctx context.Context, userID int64, request model.RefreshTokenRequest) (jwtToken string, err error)
}

type userUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

type PostUsecase interface {
	CreatePost(ctx context.Context, userID int64, req model.CreatePostRequest) (err error)
	GetPostByID(ctx context.Context, postID int64) (post model.GetPostResponse, err error)
	GetAllPost(ctx context.Context, pageSize, pageIndex int) (posts model.GetAllPostResponse, err error)
	CreateComment(ctx context.Context, postID, userID int64, request model.CreateCommentRequest) (err error)
	UpsertUserActivity(ctx context.Context, postID, userID int64, request model.UserActivityRequest) (err error)
}

type postUsecase struct {
	postRepository repository.PostRepository
}

func NewPostUsecase(postRepository repository.PostRepository) PostUsecase {
	return &postUsecase{
		postRepository: postRepository,
	}
}
