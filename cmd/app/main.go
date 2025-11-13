package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/jootiee/avito-test-2025/internal/config"
	"github.com/jootiee/avito-test-2025/internal/transport"
)

func main() {
	_ = godotenv.Load()
	config := config.NewConfig()
	config.LoadFromEnv()

	if err := transport.StartHTTPServer(config); err != nil {
		log.Fatal(err)
	}
}
