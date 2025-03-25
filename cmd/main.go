package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/suhriar/blog-mono-api/cmd/app"
	"github.com/suhriar/blog-mono-api/config"
	"github.com/suhriar/blog-mono-api/config/database"
	"github.com/suhriar/blog-mono-api/pkg/logger"
)

func main() {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch system signals for shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		oscall := <-ch
		log.Warn().Msgf("system call:%+v", oscall)
		cancel()
	}()

	// Load environment variables
	config.LoadConfig()

	// Initialize Logging
	logger.InitializeLogger(config.AppConfig)

	// Initialize DB
	db, err := database.NewMySQLConnection(config.AppConfig)
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()

	// Router setup
	router := mux.NewRouter()

	app.NewApp(router, db)

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.AppConfig.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Server in a goroutine
	go func() {
		log.Info().Msg(fmt.Sprintf("Server running on port %s", config.AppConfig.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg(fmt.Sprintf("Could not listen on %s: %v", config.AppConfig.Server.Port, err))
		}
	}()

	<-ctx.Done()

	// Graceful Shutdown
	gracefulShutdownPeriod := 30 * time.Second
	log.Warn().Msg("shutting down http server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownPeriod)
	defer cancel()

	log.Warn().Msg("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown HTTP server gracefully")
	}

	log.Info().Msg("Server shut down successfully")
}
