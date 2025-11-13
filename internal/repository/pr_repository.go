package repository

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type PRRepository interface {
	CreatePR(ctx context.Context, pr *domain.PullRequest) error
	GetPR(ctx context.Context, prID string) (*domain.PullRequest, error)
	UpdatePR(ctx context.Context, pr *domain.PullRequest) error
	PRExists(ctx context.Context, prID string) (bool, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}
