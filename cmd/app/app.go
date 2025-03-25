package app

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/suhriar/blog-mono-api/internal/delivery/rest"
	repository "github.com/suhriar/blog-mono-api/internal/repository/mysql"
	"github.com/suhriar/blog-mono-api/internal/usecase"
)

func NewApp(router *mux.Router, db *sql.DB) {
	// init repo
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// init usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	postUsecase := usecase.NewPostUsecase(postRepo)

	// init handler
	userHandler := rest.NewUserHandler(userUsecase)
	postHandler := rest.NewPostHandler(postUsecase)

	// regis rest
	rest.RegisterRoutes(router, userHandler, postHandler)
}
