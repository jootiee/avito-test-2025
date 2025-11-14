package transport

import (
	"fmt"
	"net/http"

	"github.com/jootiee/avito-test-2025/internal/config"
	"github.com/jootiee/avito-test-2025/internal/handler"
	"github.com/jootiee/avito-test-2025/internal/service"
	"github.com/jootiee/avito-test-2025/internal/storage"
	"github.com/jootiee/avito-test-2025/pkg/logger"
)

// StartHTTPServer initializes and starts the HTTP server
func StartHTTPServer(cfg *config.Config) error {
	log := logger.New(cfg.LogLevel)

	store, err := storage.New(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	log.Info("Connected to storage")

	svc := service.New(store.TeamRepo(), store.UserRepo(), store.PRRepo())

	h := handler.New(svc, log)

	log.Info("Starting server on %s", cfg.BindAddr)
	if err := http.ListenAndServe(cfg.BindAddr, h); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
