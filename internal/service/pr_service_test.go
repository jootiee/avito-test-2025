package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

func setupPRTestData(ctx context.Context, teamRepo *mockTeamRepository, userRepo *mockUserRepository) {
	// Create team
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
			{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
			{UserID: "user4", Username: "dave", TeamName: "backend", IsActive: false},
		},
	}
	teamRepo.CreateTeam(ctx, team)

	// Create users
	for _, user := range team.Members {
		userCopy := user
		userRepo.UpsertUser(ctx, &userCopy)
	}
}

func TestPRService_CreatePR_Success(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	pr, err := service.CreatePR(ctx, "pr1", "Add feature", "author1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pr.PullRequestID != "pr1" {
		t.Errorf("expected PR ID 'pr1', got %s", pr.PullRequestID)
	}

	if pr.Status != domain.PRStatusOpen {
		t.Errorf("expected status OPEN, got %s", pr.Status)
	}

	if pr.AuthorID != "author1" {
		t.Errorf("expected author 'author1', got %s", pr.AuthorID)
	}

	// Should assign up to 2 active reviewers (not the author, not inactive)
	if len(pr.AssignedReviewers) > 2 {
		t.Errorf("expected max 2 reviewers, got %d", len(pr.AssignedReviewers))
	}

	// Verify reviewers don't include author
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == "author1" {
			t.Error("author should not be assigned as reviewer")
		}
	}
}

func TestPRService_CreatePR_PRAlreadyExists(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	// Create first time
	service.CreatePR(ctx, "pr1", "Add feature", "author1")

	// Try to create again
	_, err := service.CreatePR(ctx, "pr1", "Add feature", "author1")

	if err == nil {
		t.Fatal("expected error when PR already exists")
	}

	if err.Error() != "PR already exists" {
		t.Errorf("expected 'PR already exists', got %v", err)
	}
}

func TestPRService_CreatePR_AuthorNotFound(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	_, err := service.CreatePR(ctx, "pr1", "Add feature", "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent author")
	}

	// Mock returns "user not found" which gets wrapped as "author not found" by service
	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got %v", err)
	}
}

func TestPRService_CreatePR_TeamNotFound(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	// Create user without team
	user := &domain.User{
		UserID:   "author1",
		Username: "alice",
		TeamName: "nonexistent",
		IsActive: true,
	}
	userRepo.UpsertUser(ctx, user)

	_, err := service.CreatePR(ctx, "pr1", "Add feature", "author1")

	if err == nil {
		t.Fatal("expected error for nonexistent team")
	}

	if err.Error() != "team not found" {
		t.Errorf("expected 'team not found', got %v", err)
	}
}

func TestPRService_CreatePR_NoActiveMembersToAssign(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	// Team with only author and inactive members
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: false},
		},
	}
	teamRepo.CreateTeam(ctx, team)

	for _, user := range team.Members {
		userCopy := user
		userRepo.UpsertUser(ctx, &userCopy)
	}

	pr, err := service.CreatePR(ctx, "pr1", "Add feature", "author1")

	if err != nil {
		t.Fatalf("expected no error even with no candidates, got %v", err)
	}

	// Should have no reviewers assigned
	if len(pr.AssignedReviewers) != 0 {
		t.Errorf("expected 0 reviewers, got %d", len(pr.AssignedReviewers))
	}
}

func TestPRService_MergePR_Success(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	// Create PR
	createdPR, _ := service.CreatePR(ctx, "pr1", "Add feature", "author1")

	// Merge PR
	mergedPR, err := service.MergePR(ctx, "pr1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mergedPR.Status != domain.PRStatusMerged {
		t.Errorf("expected status MERGED, got %s", mergedPR.Status)
	}

	if mergedPR.MergedAt == nil {
		t.Error("expected MergedAt to be set")
	}

	// Verify PR ID unchanged
	if mergedPR.PullRequestID != createdPR.PullRequestID {
		t.Error("PR ID should remain unchanged")
	}
}

func TestPRService_MergePR_AlreadyMerged(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	// Create and merge PR
	service.CreatePR(ctx, "pr1", "Add feature", "author1")
	firstMerge, _ := service.MergePR(ctx, "pr1")

	// Merge again (idempotent)
	secondMerge, err := service.MergePR(ctx, "pr1")

	if err != nil {
		t.Fatalf("expected no error on second merge, got %v", err)
	}

	if secondMerge.Status != domain.PRStatusMerged {
		t.Error("PR should remain merged")
	}

	// MergedAt should not change
	if !firstMerge.MergedAt.Equal(*secondMerge.MergedAt) {
		t.Error("MergedAt should not change on second merge")
	}
}

func TestPRService_MergePR_PRNotFound(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	_, err := service.MergePR(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent PR")
	}

	if err.Error() != "PR not found" {
		t.Errorf("expected 'PR not found', got %v", err)
	}
}

func TestPRService_ReassignReviewer_Success(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	// Setup team with enough members: author + 2 assigned + 1 available
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
			{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
			{UserID: "user5", Username: "eve", TeamName: "backend", IsActive: true}, // Available for reassignment
		},
	}
	teamRepo.CreateTeam(ctx, team)
	for _, user := range team.Members {
		userCopy := user
		userRepo.UpsertUser(ctx, &userCopy)
	}

	// Create PR with reviewers
	pr := domain.NewPullRequest("pr1", "Add feature", "author1")
	pr.AssignedReviewers = []string{"user2", "user3"}
	prRepo.CreatePR(ctx, pr)

	// Reassign user2
	updatedPR, newUserID, err := service.ReassignReviewer(ctx, "pr1", "user2")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newUserID == "" {
		t.Error("expected new user ID to be returned")
	}

	// Verify user2 is no longer assigned
	for _, reviewer := range updatedPR.AssignedReviewers {
		if reviewer == "user2" {
			t.Error("user2 should no longer be assigned")
		}
	}

	// Verify new user is assigned and is valid
	found := false
	for _, reviewer := range updatedPR.AssignedReviewers {
		if reviewer == newUserID {
			found = true
		}
	}
	if !found {
		t.Error("new user should be in reviewers list")
	}

	// New user should not be author or inactive
	if newUserID == "author1" {
		t.Error("author should not be assigned")
	}
	if newUserID == "user4" {
		t.Error("inactive user should not be assigned")
	}
}

func TestPRService_ReassignReviewer_PRNotFound(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	_, _, err := service.ReassignReviewer(ctx, "nonexistent", "user2")

	if err == nil {
		t.Fatal("expected error for nonexistent PR")
	}

	if err.Error() != "PR not found" {
		t.Errorf("expected 'PR not found', got %v", err)
	}
}

func TestPRService_ReassignReviewer_PRMerged(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	// Create and merge PR
	pr := domain.NewPullRequest("pr1", "Add feature", "author1")
	pr.Status = domain.PRStatusMerged
	pr.AssignedReviewers = []string{"user2", "user3"}
	prRepo.CreatePR(ctx, pr)

	_, _, err := service.ReassignReviewer(ctx, "pr1", "user2")

	if err == nil {
		t.Fatal("expected error when reassigning on merged PR")
	}

	if err.Error() != "cannot reassign on merged PR" {
		t.Errorf("expected 'cannot reassign on merged PR', got %v", err)
	}
}

func TestPRService_ReassignReviewer_ReviewerNotAssigned(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()
	setupPRTestData(ctx, teamRepo, userRepo)

	// Create PR without user2 as reviewer
	pr := domain.NewPullRequest("pr1", "Add feature", "author1")
	pr.AssignedReviewers = []string{"user3"}
	prRepo.CreatePR(ctx, pr)

	_, _, err := service.ReassignReviewer(ctx, "pr1", "user2")

	if err == nil {
		t.Fatal("expected error when user not assigned")
	}

	if err.Error() != "reviewer is not assigned to this PR" {
		t.Errorf("expected 'reviewer is not assigned to this PR', got %v", err)
	}
}

func TestPRService_ReassignReviewer_NoCandidates(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	// Small team: author + 2 reviewers only
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
			{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
		},
	}
	teamRepo.CreateTeam(ctx, team)
	for _, user := range team.Members {
		userCopy := user
		userRepo.UpsertUser(ctx, &userCopy)
	}

	// PR with both possible reviewers assigned
	pr := domain.NewPullRequest("pr1", "Add feature", "author1")
	pr.AssignedReviewers = []string{"user2", "user3"}
	prRepo.CreatePR(ctx, pr)

	_, _, err := service.ReassignReviewer(ctx, "pr1", "user2")

	if err == nil {
		t.Fatal("expected error when no candidates available")
	}

	if err.Error() != "no active replacement candidate in team" {
		t.Errorf("expected 'no active replacement candidate in team', got %v", err)
	}
}

func TestPRService_ReassignReviewer_RepositoryError(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	// Setup team with enough members for reassignment
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
			{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
			{UserID: "user5", Username: "eve", TeamName: "backend", IsActive: true},
		},
	}
	teamRepo.CreateTeam(ctx, team)
	for _, user := range team.Members {
		userCopy := user
		userRepo.UpsertUser(ctx, &userCopy)
	}

	// Create PR
	pr := domain.NewPullRequest("pr1", "Add feature", "author1")
	pr.AssignedReviewers = []string{"user2", "user3"}
	prRepo.CreatePR(ctx, pr)

	// Set error for update operation
	prRepo.updateErr = errors.New("database error")

	_, _, err := service.ReassignReviewer(ctx, "pr1", "user2")

	if err == nil {
		t.Fatal("expected repository error")
	}

	if err.Error() != "database error" {
		t.Errorf("expected 'database error', got %v", err)
	}
}

func TestPRService_GetPR_Success(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	// Create PR directly
	pr := domain.NewPullRequest("pr1", "Add feature", "author1")
	prRepo.CreatePR(ctx, pr)

	retrievedPR, err := service.GetPR(ctx, "pr1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedPR.PullRequestID != "pr1" {
		t.Errorf("expected PR ID 'pr1', got %s", retrievedPR.PullRequestID)
	}
}

func TestPRService_GetPR_NotFound(t *testing.T) {
	prRepo := newMockPRRepository()
	userRepo := newMockUserRepository()
	teamRepo := newMockTeamRepository()
	service := NewPRService(prRepo, userRepo, teamRepo)

	ctx := context.Background()

	_, err := service.GetPR(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent PR")
	}

	if err.Error() != "PR not found" {
		t.Errorf("expected 'PR not found', got %v", err)
	}
}
