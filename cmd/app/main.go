package main

import (
	"github.com/jootiee/avito-test-2025/internal/transport"
)

// @title PR Reviewer Assignment Service
// @version 1.0
// @description Service for automatic PR reviewer assignment based on team membership and workload
// @host localhost:8080
// @BasePath /
func main() {
	transport.StartHTTPServer()
}
