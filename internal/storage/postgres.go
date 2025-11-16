package storage

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/repository/postgres"
	"github.com/jootiee/avito-test-2025/internal/service"
)

// postgresStorage wraps postgres.DB and implements Storage interface
type postgresStorage struct {
	db *postgres.DB
}

// New creates a new storage instance using PostgreSQL
// The storage type is decided here - to switch to a different storage,
// only this function needs to be changed
func New(connString string) (Storage, error) {
	db, err := postgres.NewPostgresDB(connString)
	if err != nil {
		return nil, err
	}

	return &postgresStorage{db: db}, nil
}

func (s *postgresStorage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *postgresStorage) Close() error {
	return s.db.Close()
}

func (s *postgresStorage) TeamRepo() service.TeamRepository {
	return s.db.Team
}

func (s *postgresStorage) UserRepo() service.UserRepository {
	return s.db.User
}

func (s *postgresStorage) PRRepo() service.PullRequestRepository {
	return s.db.PR
}
