package transport

import (
	"fmt"
	"net/http"

	"github.com/jootiee/avito-test-2025/internal/config"
	"github.com/jootiee/avito-test-2025/internal/handler"
	"github.com/jootiee/avito-test-2025/internal/repository/postgres"
	"github.com/jootiee/avito-test-2025/internal/service"
	"github.com/jootiee/avito-test-2025/pkg/logger"
)

// StartHTTPServer initializes and starts the HTTP server
func StartHTTPServer(cfg *config.Config) error {
	log := logger.New(cfg.LogLevel)

	dbURL := cfg.DatabaseURL
	db, err := postgres.NewPostgresDB(dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	log.Info("Connected to database")

	teamService := service.NewTeamService(db.Team, db.User)
	userService := service.NewUserService(db.User, db.PR)
	prService := service.NewPRService(db.PR, db.User, db.Team)

	h := handler.New(teamService, userService, prService, log)

	log.Info("Starting server on %s", cfg.BindAddr)
	if err := http.ListenAndServe(cfg.BindAddr, h); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
