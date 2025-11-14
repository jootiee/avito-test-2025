package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jootiee/avito-test-2025/internal/domain"
)

type teamRepo struct {
	pool     *pgxpool.Pool
	userRepo *userRepo
}

func NewTeamRepository(pool *pgxpool.Pool, userRepo *userRepo) *teamRepo {
	return &teamRepo{
		pool:     pool,
		userRepo: userRepo,
	}
}

func (r *teamRepo) CreateTeam(ctx context.Context, team *domain.Team) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO teams (team_name) VALUES ($1)`, team.TeamName)
	return err
}

func (r *teamRepo) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	members, err := r.userRepo.GetTeamMembers(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		exists, err := r.TeamExists(ctx, teamName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, sql.ErrNoRows
		}
	}
	return &domain.Team{TeamName: teamName, Members: members}, nil
}

func (r *teamRepo) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, teamName).Scan(&exists)
	return exists, err
}
