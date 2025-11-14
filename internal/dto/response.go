package dto

import (
	"github.com/jootiee/avito-test-2025/internal/domain"
)

// TeamResponse wraps a team for API responses
type TeamResponse struct {
	Team *domain.Team `json:"team"`
}

// UserResponse wraps a user for API responses
type UserResponse struct {
	User *domain.User `json:"user"`
}

// PRResponse wraps a pull request for API responses
type PRResponse struct {
	PR *domain.PullRequest `json:"pr"`
}

// ReassignResponse includes PR and the new reviewer ID
type ReassignResponse struct {
	PR         *domain.PullRequest `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

// PRShort represents a shortened PR for list responses
type PRShort struct {
	PullRequestID   string          `json:"pull_request_id"`
	PullRequestName string          `json:"pull_request_name"`
	AuthorID        string          `json:"author_id"`
	Status          domain.PRStatus `json:"status"`
}

// UserReviewsResponse represents PRs where user is assigned as reviewer
type UserReviewsResponse struct {
	UserID       string    `json:"user_id"`
	PullRequests []PRShort `json:"pull_requests"`
}

// UserStats represents statistics for a single user
type UserStats struct {
	UserID          string `json:"user_id"`
	AssignmentCount int    `json:"assignment_count"`
}

// PRStats represents statistics for a single pull request
type PRStats struct {
	PullRequestID string `json:"pull_request_id"`
	ReviewerCount int    `json:"reviewer_count"`
}

// StatsResponse represents statistics about reviewer assignments
type StatsResponse struct {
	UserStats []UserStats `json:"user_stats"`
	PRStats   []PRStats   `json:"pr_stats"`
}
