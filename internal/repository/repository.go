package repository

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type Repository interface {
	// Health check
	Ping(ctx context.Context) error
	Close() error

	// Team operations
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)

	// User operations
	UpsertUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)

	// PullRequest operations
	CreatePR(ctx context.Context, pr *domain.PullRequest) error
	GetPR(ctx context.Context, prID string) (*domain.PullRequest, error)
	UpdatePR(ctx context.Context, pr *domain.PullRequest) error
	PRExists(ctx context.Context, prID string) (bool, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}
