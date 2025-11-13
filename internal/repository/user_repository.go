package repository

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type UserRepository interface {
	UpsertUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)
}
