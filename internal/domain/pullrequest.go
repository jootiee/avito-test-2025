package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            PRStatus   `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func NewPullRequest(id, name, authorID string) *PullRequest {
	return &PullRequest{
		PullRequestID:     id,
		PullRequestName:   name,
		AuthorID:          authorID,
		Status:            PRStatusOpen,
		AssignedReviewers: []string{},
	}
}
