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

func TestHandler_PREndpoints(t *testing.T) {
	// Setup
	teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
	userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
	prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

	// Setup team and users
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
			{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
		},
	}
	teamRepo.CreateTeam(nil, team)
	for _, user := range team.Members {
		userCopy := user
		userRepo.UpsertUser(nil, &userCopy)
	}

	svc := service.New(teamRepo, userRepo, prRepo)
	log := &mockLogger{}
	h := New(svc, log)

	// Test 1: Create PR
	t.Run("CreatePR_Success", func(t *testing.T) {
		reqBody := dto.CreatePRRequest{
			PullRequestID:   "pr1",
			PullRequestName: "Add feature",
			AuthorID:        "author1",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		var resp dto.PRResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.PR.PullRequestID != "pr1" {
			t.Errorf("expected PR ID 'pr1', got %s", resp.PR.PullRequestID)
		}
		if resp.PR.Status != domain.PRStatusOpen {
			t.Errorf("expected status OPEN, got %s", resp.PR.Status)
		}
	})

	// Test 2: Create duplicate PR
	t.Run("CreatePR_AlreadyExists", func(t *testing.T) {
		reqBody := dto.CreatePRRequest{
			PullRequestID:   "pr1",
			PullRequestName: "Add feature",
			AuthorID:        "author1",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status 409, got %d", w.Code)
		}
	})

	// Test 3: Merge PR
	t.Run("MergePR_Success", func(t *testing.T) {
		reqBody := dto.MergeRequest{
			PullRequestID: "pr1",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp dto.PRResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.PR.Status != domain.PRStatusMerged {
			t.Errorf("expected status MERGED, got %s", resp.PR.Status)
		}
	})

	// Test 4: Reassign reviewer
	t.Run("ReassignReviewer_CannotReassignMergedPR", func(t *testing.T) {
		reqBody := dto.ReassignRequest{
			PullRequestID: "pr1",
			OldUserID:     "user2",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status 409, got %d", w.Code)
		}

		var resp dto.APIError
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.Error.Code != dto.ErrCodePRMerged {
			t.Errorf("expected error code %s, got %s", dto.ErrCodePRMerged, resp.Error.Code)
		}
	})

	// Test 5: Create PR - author not found
	t.Run("CreatePR_AuthorNotFound", func(t *testing.T) {
		reqBody := dto.CreatePRRequest{
			PullRequestID:   "pr2",
			PullRequestName: "Another feature",
			AuthorID:        "nonexistent",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	// Test 6: Merge PR - not found
	t.Run("MergePR_NotFound", func(t *testing.T) {
		reqBody := dto.MergeRequest{
			PullRequestID: "nonexistent",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	// Test 7: Invalid JSON
	t.Run("CreatePR_InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
