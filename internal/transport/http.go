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

	// Initialize storage - the specific implementation (postgres) is decided in storage package
	store, err := storage.New(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	log.Info("Connected to database")

	teamService := service.NewTeamService(store.TeamRepo(), store.UserRepo())
	userService := service.NewUserService(store.UserRepo(), store.PRRepo())
	prService := service.NewPRService(store.PRRepo(), store.UserRepo(), store.TeamRepo())

	h := handler.New(teamService, userService, prService, log)

	log.Info("Starting server on %s", cfg.BindAddr)
	if err := http.ListenAndServe(cfg.BindAddr, h); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
