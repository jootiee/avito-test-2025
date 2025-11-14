package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
)

func TestHandler_UserEndpoints(t *testing.T) {
	// Setup
	teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
	userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
	prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

	// Add test user
	userRepo.UpsertUser(nil, &domain.User{
		UserID:   "user1",
		Username: "alice",
		TeamName: "backend",
		IsActive: true,
	})

	svc := service.New(teamRepo, userRepo, prRepo)
	log := &mockLogger{}
	h := New(svc, log)

	// Test 1: Set user active
	t.Run("SetUserActive_Success", func(t *testing.T) {
		reqBody := dto.SetActiveRequest{
			UserID:   "user1",
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp dto.UserResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.User.IsActive {
			t.Error("expected user to be inactive")
		}
	})

	// Test 2: Set user active - not found
	t.Run("SetUserActive_NotFound", func(t *testing.T) {
		reqBody := dto.SetActiveRequest{
			UserID:   "nonexistent",
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	// Test 3: Get user reviews
	t.Run("GetUserReviews_Success", func(t *testing.T) {
		// Add PR with user as reviewer
		pr := domain.NewPullRequest("pr1", "Feature", "author1")
		pr.AssignedReviewers = []string{"user1"}
		prRepo.CreatePR(nil, pr)

		req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=user1", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp dto.UserReviewsResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.PullRequests) != 1 {
			t.Errorf("expected 1 PR, got %d", len(resp.PullRequests))
		}
	})

	// Test 4: Get user reviews - missing user_id
	t.Run("GetUserReviews_MissingUserID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	// Test 5: Invalid JSON
	t.Run("SetUserActive_InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
