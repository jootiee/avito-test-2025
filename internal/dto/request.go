package dto

// TeamAddRequest represents the request to create a team
type TeamAddRequest struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

// TeamMemberDTO represents a team member in requests
type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// SetActiveRequest represents the request to set user active status
type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// CreatePullRequestRequest represents the request to create a pull request
type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

// MergeRequest represents the request to merge a PR
type MergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

// ReassignRequest represents the request to reassign a reviewer
type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}
