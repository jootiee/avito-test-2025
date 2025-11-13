package transport

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/jootiee/avito-test-2025/internal/config"
	"github.com/jootiee/avito-test-2025/internal/handler"
	"github.com/jootiee/avito-test-2025/internal/repository/postgres"
	"github.com/jootiee/avito-test-2025/internal/service"
)

// Initializes and starts the HTTP server
func StartHTTPServer(cfg *config.Config) error {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	dbURL := cfg.DatabaseURL
	db, err := postgres.NewPostgresDB(dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	teamService := service.NewTeamService(db.Team, db.User)
	userService := service.NewUserService(db.User, db.PR)
	prService := service.NewPRService(db.PR, db.User, db.Team)

	h := handler.NewHandler(teamService, userService, prService, logger)

	log.Printf("Starting server on %s", cfg.BindAddr)
	if err := http.ListenAndServe(cfg.BindAddr, h); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
