package repository

import (
	"context"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}
