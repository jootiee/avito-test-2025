package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
)

// Helper to create test handler with mocked team service
func newTestHandlerWithTeamMock(
	createFunc func(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error),
	getFunc func(ctx context.Context, teamName string) (*domain.Team, error),
) *Handler {
	teamService := &service.TeamService{}
	// Note: We can't easily mock concrete services without refactoring to interfaces at handler level
	// For now, create minimal service that will work with simple tests
	svc := &service.Service{
		Team: teamService,
	}
	return New(svc, &mockLogger{})
}

func TestHandler_AddTeam_InvalidJSON(t *testing.T) {
	// Setup - create handler with actual service but we'll fail before calling it
	svc := &service.Service{}
	log := &mockLogger{}
	h := New(svc, log)

	// Invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp dto.APIError
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error.Code != dto.ErrCodeNotFound {
		t.Errorf("expected error code %s, got %s", dto.ErrCodeNotFound, resp.Error.Code)
	}
}

func TestHandler_GetTeam_MissingTeamName(t *testing.T) {
	// Setup
	svc := &service.Service{}
	log := &mockLogger{}
	h := New(svc, log)

	// Create request without team_name query param
	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	w := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp dto.APIError
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error.Message != "team_name required" {
		t.Errorf("expected 'team_name required' message, got %s", resp.Error.Message)
	}
}

// Integration-style test with real service but mocked repositories
func TestHandler_TeamEndpoints_WithMockedRepos(t *testing.T) {
	// Create mocked repositories
	teamRepo := &mockTeamRepo{
		teams: make(map[string]*domain.Team),
	}
	userRepo := &mockUserRepo{
		users: make(map[string]*domain.User),
	}
	prRepo := &mockPRRepo{
		prs: make(map[string]*domain.PullRequest),
	}

	// Create real services with mocked repositories
	svc := service.New(teamRepo, userRepo, prRepo)
	log := &mockLogger{}
	h := New(svc, log)

	// Test 1: Add team successfully
	t.Run("AddTeam", func(t *testing.T) {
		reqBody := dto.TeamAddRequest{
			TeamName: "backend",
			Members: []dto.TeamMemberDTO{
				{UserID: "user1", Username: "alice", IsActive: true},
				{UserID: "user2", Username: "bob", IsActive: true},
			},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		var resp dto.TeamResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.Team.TeamName != "backend" {
			t.Errorf("expected team name 'backend', got %s", resp.Team.TeamName)
		}
	})

	// Test 2: Get team
	t.Run("GetTeam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=backend", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	// Test 3: Get non-existent team
	t.Run("GetTeam_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=nonexistent", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	// Test 4: Add duplicate team
	t.Run("AddTeam_AlreadyExists", func(t *testing.T) {
		reqBody := dto.TeamAddRequest{
			TeamName: "backend",
			Members:  []dto.TeamMemberDTO{{UserID: "user1", Username: "alice", IsActive: true}},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// Simple mock repository implementations for handler testing
type mockTeamRepo struct {
	teams map[string]*domain.Team
}

func (m *mockTeamRepo) CreateTeam(ctx context.Context, team *domain.Team) error {
	m.teams[team.TeamName] = team
	return nil
}

func (m *mockTeamRepo) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	team, exists := m.teams[teamName]
	if !exists {
		return nil, errors.New("team not found")
	}
	return team, nil
}

func (m *mockTeamRepo) TeamExists(ctx context.Context, teamName string) (bool, error) {
	_, exists := m.teams[teamName]
	return exists, nil
}

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) UpsertUser(ctx context.Context, user *domain.User) error {
	m.users[user.UserID] = user
	return nil
}

func (m *mockUserRepo) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepo) SetUserActive(ctx context.Context, userID string, isActive bool) error {
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.IsActive = isActive
	return nil
}

func (m *mockUserRepo) GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	var members []domain.User
	for _, user := range m.users {
		if user.TeamName == teamName {
			members = append(members, *user)
		}
	}
	return members, nil
}

type mockPRRepo struct {
	prs map[string]*domain.PullRequest
}

func (m *mockPRRepo) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	m.prs[pr.PullRequestID] = pr
	return nil
}

func (m *mockPRRepo) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, exists := m.prs[prID]
	if !exists {
		return nil, errors.New("PR not found")
	}
	return pr, nil
}

func (m *mockPRRepo) UpdatePR(ctx context.Context, pr *domain.PullRequest) error {
	m.prs[pr.PullRequestID] = pr
	return nil
}

func (m *mockPRRepo) PRExists(ctx context.Context, prID string) (bool, error) {
	_, exists := m.prs[prID]
	return exists, nil
}

func (m *mockPRRepo) GetPRsByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	var prs []*domain.PullRequest
	for _, pr := range m.prs {
		for _, reviewer := range pr.AssignedReviewers {
			if reviewer == userID {
				prs = append(prs, pr)
				break
			}
		}
	}
	return prs, nil
}
