package repository

import (
	"context"
)

// Repository defines the interface for database health checks
type Repository interface {
	Ping(ctx context.Context) error
	Close() error
}
