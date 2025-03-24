package mysql

import (
	"context"
	"database/sql"

	"github.com/suhriar/blog-mono-api/model"
)

type UserRepository interface {
	GetUser(ctx context.Context, email, username string, userID int64) (user model.User, err error)
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
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}
