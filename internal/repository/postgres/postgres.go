package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool

	Team *TeamRepo
	User *UserRepo
	PR   *PRRepo
}

// NewPostgresDB creates a new PostgreSQL connection pool and initializes all repositories
func NewPostgresDB(connString string) (*DB, error) {
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

	return &DB{
		pool: pool,
		Team: teamRepo,
		User: userRepo,
		PR:   prRepo,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}
