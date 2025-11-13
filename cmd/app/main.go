package main

import (
	"github.com/jootiee/avito-test-2025/pkg/httpserver"
)

func main() {
	if err := httpserver.Start(); err != nil {
		return
	}
}