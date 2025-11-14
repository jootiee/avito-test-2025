package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jootiee/avito-test-2025/internal/domain"
)

type prRepo struct {
	pool *pgxpool.Pool
}

func NewPRRepository(pool *pgxpool.Pool) *prRepo {
	return &prRepo{pool: pool}
}

func (r *prRepo) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, reviewersJSON, pr.CreatedAt)
	return err
}

func (r *prRepo) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	var reviewersJSON []byte
	err := r.pool.QueryRow(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_requests WHERE pull_request_id = $1
	`, prID).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(reviewersJSON, &pr.AssignedReviewers); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *prRepo) UpdatePR(ctx context.Context, pr *domain.PullRequest) error {
	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE pull_requests SET
			pull_request_name = $2,
			author_id = $3,
			status = $4,
			assigned_reviewers = $5,
			merged_at = $6
		WHERE pull_request_id = $1
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, reviewersJSON, pr.MergedAt)
	return err
}

func (r *prRepo) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`, prID).Scan(&exists)
	return exists, err
}

func (r *prRepo) GetPRsByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_requests
		WHERE assigned_reviewers::jsonb ? $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		var reviewersJSON []byte
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(reviewersJSON, &pr.AssignedReviewers); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, rows.Err()
}

func (r *prRepo) GetUserAssignmentCounts(ctx context.Context) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT reviewer_id, COUNT(*) as assignment_count
		FROM pull_requests,
		jsonb_array_elements_text(assigned_reviewers) AS reviewer_id
		GROUP BY reviewer_id
		ORDER BY assignment_count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var userID string
		var count int
		if err := rows.Scan(&userID, &count); err != nil {
			return nil, err
		}
		counts[userID] = count
	}
	return counts, rows.Err()
}

func (r *prRepo) GetPRReviewerCounts(ctx context.Context) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT pull_request_id, jsonb_array_length(assigned_reviewers) as reviewer_count
		FROM pull_requests
		ORDER BY reviewer_count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var prID string
		var count int
		if err := rows.Scan(&prID, &count); err != nil {
			return nil, err
		}
		counts[prID] = count
	}
	return counts, rows.Err()
}
