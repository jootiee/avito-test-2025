package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jootiee/avito-test-2025/internal/domain"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

func (r *userRepo) UpsertUser(ctx context.Context, user *domain.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
	`, user.UserID, user.Username, user.TeamName, user.IsActive)
	return err
}

func (r *userRepo) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users WHERE user_id = $1
	`, userID).Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) SetUserActive(ctx context.Context, userID string, isActive bool) error {
	res, err := r.pool.Exec(ctx, `UPDATE users SET is_active = $1, updated_at = NOW() WHERE user_id = $2`, isActive, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepo) GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	return members, rows.Err()
}
