package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

func TestUserService_SetUserActive_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	// Setup existing user
	user := &domain.User{
		UserID:   "user1",
		Username: "alice",
		TeamName: "backend",
		IsActive: true,
	}
	userRepo.UpsertUser(ctx, user)

	// Set user inactive
	updatedUser, err := service.SetUserActive(ctx, "user1", false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedUser.IsActive {
		t.Error("expected user to be inactive")
	}

	// Verify in repository
	storedUser, _ := userRepo.GetUser(ctx, "user1")
	if storedUser.IsActive {
		t.Error("expected stored user to be inactive")
	}
}

func TestUserService_SetUserActive_UserNotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	_, err := service.SetUserActive(ctx, "nonexistent", false)

	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got %v", err)
	}
}

func TestUserService_SetUserActive_RepositoryError(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	// Setup existing user
	user := &domain.User{
		UserID:   "user1",
		Username: "alice",
		TeamName: "backend",
		IsActive: true,
	}
	userRepo.UpsertUser(ctx, user)

	// Inject error
	userRepo.setActiveErr = errors.New("database error")

	_, err := service.SetUserActive(ctx, "user1", false)

	if err == nil {
		t.Fatal("expected repository error")
	}

	if err.Error() != "database error" {
		t.Errorf("expected 'database error', got %v", err)
	}
}

func TestUserService_GetUser_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	// Setup existing user
	user := &domain.User{
		UserID:   "user1",
		Username: "alice",
		TeamName: "backend",
		IsActive: true,
	}
	userRepo.UpsertUser(ctx, user)

	retrievedUser, err := service.GetUser(ctx, "user1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedUser.UserID != "user1" {
		t.Errorf("expected user ID 'user1', got %s", retrievedUser.UserID)
	}

	if retrievedUser.Username != "alice" {
		t.Errorf("expected username 'alice', got %s", retrievedUser.Username)
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	_, err := service.GetUser(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got %v", err)
	}
}

func TestUserService_GetUserReviews_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	// Setup PRs
	pr1 := &domain.PullRequest{
		PullRequestID:     "pr1",
		PullRequestName:   "Feature A",
		AuthorID:          "author1",
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"user1", "user2"},
	}
	pr2 := &domain.PullRequest{
		PullRequestID:     "pr2",
		PullRequestName:   "Feature B",
		AuthorID:          "author2",
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"user1"},
	}
	pr3 := &domain.PullRequest{
		PullRequestID:     "pr3",
		PullRequestName:   "Feature C",
		AuthorID:          "author1",
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"user3"},
	}

	prRepo.CreatePR(ctx, pr1)
	prRepo.CreatePR(ctx, pr2)
	prRepo.CreatePR(ctx, pr3)

	// Get reviews for user1
	prs, err := service.GetUserReviews(ctx, "user1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(prs) != 2 {
		t.Errorf("expected 2 PRs for user1, got %d", len(prs))
	}

	// Verify correct PRs are returned
	prIDs := make(map[string]bool)
	for _, pr := range prs {
		prIDs[pr.PullRequestID] = true
	}

	if !prIDs["pr1"] || !prIDs["pr2"] {
		t.Error("expected pr1 and pr2 to be returned")
	}
}

func TestUserService_GetUserReviews_Empty(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	prs, err := service.GetUserReviews(ctx, "user1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(prs) != 0 {
		t.Errorf("expected 0 PRs, got %d", len(prs))
	}
}

func TestUserService_GetUserReviews_RepositoryError(t *testing.T) {
	userRepo := newMockUserRepository()
	prRepo := newMockPRRepository()
	prRepo.getErr = errors.New("database error")
	service := NewUserService(userRepo, prRepo)

	ctx := context.Background()

	_, err := service.GetUserReviews(ctx, "user1")

	if err == nil {
		t.Fatal("expected repository error")
	}

	if err.Error() != "database error" {
		t.Errorf("expected 'database error', got %v", err)
	}
}
