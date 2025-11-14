package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, initiating graceful shutdown...")
		cancel()
	}()

	if err := transport.StartHTTPServer(ctx, config); err != nil {
		log.Fatal(err)
	}
}
