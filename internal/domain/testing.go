package domain

import "testing"

// TestTeam creates a test team for testing purposes
func TestTeam(t *testing.T) *Team {
	t.Helper()

	return &Team{
		TeamName: "backend",
		Members: []User{
			{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
			{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
			{UserID: "u3", Username: "Charlie", TeamName: "backend", IsActive: true},
		},
	}
}

// TestUser creates a test user for testing purposes
func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}
}

// TestPullRequest creates a test PR for testing purposes
func TestPullRequest(t *testing.T) *PullRequest {
	t.Helper()

	return NewPullRequest("pr1", "Add feature", "u1")
}
