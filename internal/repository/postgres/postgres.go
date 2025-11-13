package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jootiee/avito-test-2025/internal/repository"
)

type PostgresDB struct {
	pool *pgxpool.Pool

	Team repository.TeamRepository
	User repository.UserRepository
	PR   repository.PRRepository
}

// NewPostgresDB creates a new PostgreSQL connection pool and initializes all repositories
func NewPostgresDB(connString string) (*PostgresDB, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize repositories
	userRepo := NewUserRepository(pool)
	teamRepo := NewTeamRepository(pool, userRepo)
	prRepo := NewPRRepository(pool)

	return &PostgresDB{
		pool: pool,
		Team: teamRepo,
		User: userRepo,
		PR:   prRepo,
	}, nil
}

func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *PostgresDB) Close() error {
	db.pool.Close()
	return nil
}
