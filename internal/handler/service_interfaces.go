package handler

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

// PRService defines the interface for pull request operations
type PRService interface {
	CreatePR(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error)
	GetPR(ctx context.Context, prID string) (*domain.PullRequest, error)
}

// TeamService defines the interface for team operations
type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

// UserService defines the interface for user operations
type UserService interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}
