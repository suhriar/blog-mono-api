package usecase

import (
	"context"

	repository "github.com/suhriar/blog-mono-api/internal/repository/mysql"
	"github.com/suhriar/blog-mono-api/model"
)

type UserUsecase interface {
	SignUp(ctx context.Context, req model.SignUpRequest) (err error)
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
}

type postUsecase struct {
	postRepository repository.PostRepository
}

func NewPostUsecase(postRepository repository.PostRepository) PostUsecase {
	return &postUsecase{
		postRepository: postRepository,
	}
}
