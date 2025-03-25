package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suhriar/blog-mono-api/internal/delivery/middleware"
)

func RegisterRoutes(router *mux.Router, userHandler *UserHandler, postHandler *PostHandler) {
	router.Use(middleware.LoggingMiddleware)

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/health", HealthCheck).Methods("GET")

	jwtMiddleware := middleware.NewJWTMiddleware()

	// Register user routes
	registerUserRoutes(apiRouter, userHandler, jwtMiddleware)
	registerPostRoutes(apiRouter, postHandler, jwtMiddleware)
}

func registerUserRoutes(router *mux.Router, handler *UserHandler, jwtMiddleware *middleware.JWTMiddleware) {
	userRouter := router.PathPrefix("/users").Subrouter()

	// Public routes
	userRouter.HandleFunc("/sign-up", handler.SignUp).Methods("POST")
	userRouter.HandleFunc("/login", handler.Login).Methods("POST")

	// Protected routes
	protected := userRouter.PathPrefix("").Subrouter()
	protected.Use(jwtMiddleware.RequireAuth)
	protected.HandleFunc("/refresh", handler.Refresh).Methods("POST")
}

func registerPostRoutes(router *mux.Router, handler *PostHandler, jwtMiddleware *middleware.JWTMiddleware) {
	userRouter := router.PathPrefix("/posts").Subrouter()

	// Protected routes
	protected := userRouter.PathPrefix("").Subrouter()
	protected.Use(jwtMiddleware.RequireAuth)
	protected.HandleFunc("/create", handler.CreatePost).Methods("POST")
	protected.HandleFunc("/", handler.GetAllPost).Methods("GET")
	protected.HandleFunc("/{id:[0-9]+}", handler.GetPostByID).Methods("GET")
	protected.HandleFunc("/{id:[0-9]+}/comment", handler.CreateComment).Methods("POST")
	protected.HandleFunc("/{id:[0-9]+}/user-activity", handler.UpsertUserActivity).Methods("PUT")
}

// HealthCheck handler for the health endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
