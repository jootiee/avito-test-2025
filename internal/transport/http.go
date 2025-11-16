package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/jootiee/avito-test-2025/internal/config"
	"github.com/jootiee/avito-test-2025/internal/handler"
	"github.com/jootiee/avito-test-2025/internal/service"
	"github.com/jootiee/avito-test-2025/internal/storage"
	"github.com/jootiee/avito-test-2025/pkg/logger"
)

// StartHTTPServer initializes and starts the HTTP server with graceful shutdown
func StartHTTPServer() {
	_ = godotenv.Load()
	cfg := config.NewConfig()
	cfg.LoadFromEnv()
	log := logger.New(cfg.LogLevel, cfg.LogFormat)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Shutdown signal received, initiating graceful shutdown...")
		cancel()
	}()

	store, err := storage.New(cfg.DatabaseURL)
	if err != nil {
		log.Error("Failed to initialize storage", "error", err.Error())
		return
	}

	log.Info("Connected to storage")

	svc := service.New(store.TeamRepo(), store.UserRepo(), store.PRRepo())

	h := handler.New(svc, log)

	srv := &http.Server{
		Addr:              cfg.BindAddr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("Starting server", "address", cfg.BindAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErr:
		log.Error("Server error occurred", "error", err.Error())
		return
	case <-ctx.Done():
		log.Info("Shutting down server gracefully...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("Server shutdown error", "error", err)
		}

		log.Info("Closing database connections...")
		if err := store.Close(); err != nil {
			log.Error("Database close error", "error", err)
		}

		log.Info("Server stopped gracefully")
		return
	}
}
