package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

// TeamService handles team-related business logic
type TeamService struct {
	teamRepo TeamRepository
	userRepo UserRepository
}

// NewTeamService creates a new team service
func NewTeamService(
	teamRepo TeamRepository,
	userRepo UserRepository,
) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// Create creates a new team and upserts its members
func (s *TeamService) Create(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error) {
	exists, err := s.teamRepo.TeamExists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("team already exists")
	}

	team := domain.NewTeam(teamName, members)
	if err := s.teamRepo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}

	for i := range members {
		members[i].TeamName = teamName
		if err := s.userRepo.UpsertUser(ctx, &members[i]); err != nil {
			return nil, err
		}
	}

	return team, nil
}

// Get retrieves a team by name
func (s *TeamService) Get(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	return team, nil
}
