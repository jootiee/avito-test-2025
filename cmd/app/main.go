package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/jootiee/avito-test-2025/internal/config"
	"github.com/jootiee/avito-test-2025/internal/transport"
)

// @title PR Reviewer Assignment Service
// @version 1.0
// @description Service for automatic PR reviewer assignment based on team membership and workload
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()
	config := config.NewConfig()
	config.LoadFromEnv()

	if err := transport.StartHTTPServer(config); err != nil {
		log.Fatal(err)
	}
}
