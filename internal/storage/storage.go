package storage

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/service"
)

// Storage defines the interface for data storage operations
type Storage interface {
	// Health check
	Ping(ctx context.Context) error
	Close() error

	// Repository accessors - return interfaces from service package
	TeamRepo() service.TeamRepository
	UserRepo() service.UserRepository
	PRRepo() service.PullRequestRepository
}
