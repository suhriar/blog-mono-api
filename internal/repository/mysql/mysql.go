package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/suhriar/blog-mono-api/model"
)

type UserRepository interface {
	GetUser(ctx context.Context, email, username string, userID int64) (user model.User, err error)
	CreateUser(ctx context.Context, model model.User) (lastInsertID int64, err error)
	InsertRefreshToken(ctx context.Context, model model.RefreshToken) (lastInsertID int64, err error)
	GetRefreshToken(ctx context.Context, userID int64, now time.Time) (resp model.RefreshToken, err error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type PostRepository interface {
	CreatePost(ctx context.Context, model model.Post) (lastInsertID int64, err error)
	GetAllPost(ctx context.Context, limit, offset int) (resp model.GetAllPostResponse, err error)
	GetPostByID(ctx context.Context, id int64) (resp model.PostDetail, err error)
	CreateComment(ctx context.Context, model model.Comment) (lastInsertID int64, err error)
	GetCommentsByPostID(ctx context.Context, postID int64) (comments []model.CommentResponse, err error)
	GetUserActivity(ctx context.Context, model model.UserActivity) (resp model.UserActivity, err error)
	CreateUserActivity(ctx context.Context, model model.UserActivity) (lastInsertID int64, err error)
	UpdateUserActivity(ctx context.Context, req model.UserActivity) (err error)
	CountLikeByPostID(ctx context.Context, postID int64) (count int, err error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}
