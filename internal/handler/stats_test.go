package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
)

func TestHandler_GetStats_Success(t *testing.T) {
	// Setup
	teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
	userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
	prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

	// Create test data: 3 PRs with different reviewer assignments
	pr1 := domain.NewPullRequest("pr1", "Feature A", "author1")
	pr1.AssignedReviewers = []string{"user2", "user3"}
	prRepo.CreatePR(nil, pr1)

	pr2 := domain.NewPullRequest("pr2", "Feature B", "author2")
	pr2.AssignedReviewers = []string{"user2", "user4"}
	prRepo.CreatePR(nil, pr2)

	pr3 := domain.NewPullRequest("pr3", "Feature C", "author3")
	pr3.AssignedReviewers = []string{"user3"}
	prRepo.CreatePR(nil, pr3)

	svc := service.New(teamRepo, userRepo, prRepo)
	log := &mockLogger{}
	h := New(svc, log)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp dto.StatsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify user stats
	if len(resp.UserStats) != 3 {
		t.Errorf("expected 3 users in stats, got %d", len(resp.UserStats))
	}

	// Check that we have stats for all users
	userCountMap := make(map[string]int)
	for _, userStat := range resp.UserStats {
		userCountMap[userStat.UserID] = userStat.AssignmentCount
	}

	if userCountMap["user2"] != 2 {
		t.Errorf("expected user2 to have 2 assignments, got %d", userCountMap["user2"])
	}

	if userCountMap["user3"] != 2 {
		t.Errorf("expected user3 to have 2 assignments, got %d", userCountMap["user3"])
	}

	if userCountMap["user4"] != 1 {
		t.Errorf("expected user4 to have 1 assignment, got %d", userCountMap["user4"])
	}

	// Verify PR stats
	if len(resp.PRStats) != 3 {
		t.Errorf("expected 3 PRs in stats, got %d", len(resp.PRStats))
	}

	prCountMap := make(map[string]int)
	for _, prStat := range resp.PRStats {
		prCountMap[prStat.PullRequestID] = prStat.ReviewerCount
	}

	if prCountMap["pr1"] != 2 {
		t.Errorf("expected pr1 to have 2 reviewers, got %d", prCountMap["pr1"])
	}

	if prCountMap["pr2"] != 2 {
		t.Errorf("expected pr2 to have 2 reviewers, got %d", prCountMap["pr2"])
	}

	if prCountMap["pr3"] != 1 {
		t.Errorf("expected pr3 to have 1 reviewer, got %d", prCountMap["pr3"])
	}
}

func TestHandler_GetStats_EmptyData(t *testing.T) {
	// Setup with no PRs
	teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
	userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
	prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

	svc := service.New(teamRepo, userRepo, prRepo)
	log := &mockLogger{}
	h := New(svc, log)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp dto.StatsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.UserStats) != 0 {
		t.Errorf("expected 0 users in stats, got %d", len(resp.UserStats))
	}

	if len(resp.PRStats) != 0 {
		t.Errorf("expected 0 PRs in stats, got %d", len(resp.PRStats))
	}
}

func TestHandler_GetStats_SinglePRWithNoReviewers(t *testing.T) {
	// Setup with a PR that has no reviewers
	teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
	userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
	prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

	pr1 := domain.NewPullRequest("pr1", "Feature A", "author1")
	pr1.AssignedReviewers = []string{} // No reviewers
	prRepo.CreatePR(nil, pr1)

	svc := service.New(teamRepo, userRepo, prRepo)
	log := &mockLogger{}
	h := New(svc, log)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp dto.StatsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// No user stats since no reviewers assigned
	if len(resp.UserStats) != 0 {
		t.Errorf("expected 0 users in stats, got %d", len(resp.UserStats))
	}

	// Should have one PR with 0 reviewers
	if len(resp.PRStats) != 1 {
		t.Errorf("expected 1 PR in stats, got %d", len(resp.PRStats))
	}

	if resp.PRStats[0].ReviewerCount != 0 {
		t.Errorf("expected pr1 to have 0 reviewers, got %d", resp.PRStats[0].ReviewerCount)
	}
}
