package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
)

type TeamHandlerTestSuite struct {
	suite.Suite
}

func (s *TeamHandlerTestSuite) TestTeamEndpoints() {
	type args struct {
		method string
		path   string
		body   interface{}
	}
	type setupResult struct {
		handler *Handler
	}

	tests := []struct {
		name       string
		args       args
		setup      func() setupResult
		assertFunc func(w *httptest.ResponseRecorder)
	}{
		{
			name:  "add team invalid JSON",
			args:  args{method: http.MethodPost, path: "/team/add", body: []byte("invalid json")},
			setup: func() setupResult { return setupResult{handler: New(&service.Service{}, &mockLogger{})} },
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusBadRequest, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), dto.ErrCodeNotFound, resp.Error.Code)
			},
		},
		{
			name:  "get team missing query param",
			args:  args{method: http.MethodGet, path: "/team/get"},
			setup: func() setupResult { return setupResult{handler: New(&service.Service{}, &mockLogger{})} },
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusBadRequest, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "team_name required", resp.Error.Message)
			},
		},
		{
			name: "add team success",
			args: args{method: http.MethodPost, path: "/team/add", body: dto.TeamAddRequest{TeamName: "backend", Members: []dto.TeamMemberDTO{{UserID: "user1", Username: "alice", IsActive: true}, {UserID: "user2", Username: "bob", IsActive: true}}}},
			setup: func() setupResult {
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				h := New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
				return setupResult{handler: h}
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusCreated, w.Code)
				var resp dto.TeamResponse
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "backend", resp.Team.TeamName)
				assert.Len(s.T(), resp.Team.Members, 2)
			},
		},
		{
			name: "add team duplicate",
			args: args{method: http.MethodPost, path: "/team/add", body: dto.TeamAddRequest{TeamName: "backend", Members: []dto.TeamMemberDTO{{UserID: "user1", Username: "alice", IsActive: true}}}},
			setup: func() setupResult {
				// pre-create team
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": domain.NewTeam("backend", []domain.User{})}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				h := New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
				return setupResult{handler: h}
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusBadRequest, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), dto.ErrCodeTeamExists, resp.Error.Code)
			},
		},
		{
			name: "get team not found",
			args: args{method: http.MethodGet, path: "/team/get?team_name=missing"},
			setup: func() setupResult {
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				h := New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
				return setupResult{handler: h}
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusNotFound, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "team not found", resp.Error.Message)
			},
		},
		{
			name: "get team success",
			args: args{method: http.MethodGet, path: "/team/get?team_name=backend"},
			setup: func() setupResult {
				team := domain.NewTeam("backend", []domain.User{{UserID: "u1", Username: "alice", TeamName: "backend", IsActive: true}})
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": team}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{"u1": {UserID: "u1", Username: "alice", TeamName: "backend", IsActive: true}}}
				h := New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
				return setupResult{handler: h}
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusOK, w.Code)
				var resp domain.Team
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "backend", resp.TeamName)
				assert.Len(s.T(), resp.Members, 1)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			sr := tt.setup()
			var bodyBytes []byte
			switch b := tt.args.body.(type) {
			case nil:
			case []byte:
				bodyBytes = b
			default:
				marshaled, err := json.Marshal(b)
				s.NoError(err)
				bodyBytes = marshaled
			}
			req := httptest.NewRequest(tt.args.method, tt.args.path, bytes.NewReader(bodyBytes))
			if bodyBytes != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			sr.handler.ServeHTTP(w, req)
			tt.assertFunc(w)
		})
	}
}

func TestTeamHandler(t *testing.T) {
	suite.Run(t, new(TeamHandlerTestSuite))
}

// Simple mock repository implementations for handler testing
type mockTeamRepo struct {
	teams map[string]*domain.Team
}

func (m *mockTeamRepo) CreateTeam(_ context.Context, team *domain.Team) error {
	m.teams[team.TeamName] = team
	return nil
}

func (m *mockTeamRepo) GetTeam(_ context.Context, teamName string) (*domain.Team, error) {
	team, exists := m.teams[teamName]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return team, nil
}

func (m *mockTeamRepo) TeamExists(_ context.Context, teamName string) (bool, error) {
	_, exists := m.teams[teamName]
	return exists, nil
}

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) UpsertUser(_ context.Context, user *domain.User) error {
	m.users[user.UserID] = user
	return nil
}

func (m *mockUserRepo) GetUser(_ context.Context, userID string) (*domain.User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepo) SetUserActive(_ context.Context, userID string, isActive bool) error {
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

type mockLogger struct{}

func (m *mockLogger) Debug(_ interface{}, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})       {}
func (m *mockLogger) Warn(_ string, _ ...interface{})       {}
func (m *mockLogger) Error(_ interface{}, _ ...interface{}) {}
func (m *mockLogger) Fatal(_ interface{}, _ ...interface{}) {}

func (m *mockPRRepo) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	m.prs[pr.PullRequestID] = pr
	return nil
}

func (m *mockPRRepo) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, exists := m.prs[prID]
	if !exists {
		return nil, sql.ErrNoRows
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

func (m *mockPRRepo) GetUserAssignmentCounts(ctx context.Context) (map[string]int, error) {
	counts := make(map[string]int)
	for _, pr := range m.prs {
		for _, reviewer := range pr.AssignedReviewers {
			counts[reviewer]++
		}
	}
	return counts, nil
}

func (m *mockPRRepo) GetPRReviewerCounts(ctx context.Context) (map[string]int, error) {
	counts := make(map[string]int)
	for prID, pr := range m.prs {
		counts[prID] = len(pr.AssignedReviewers)
	}
	return counts, nil
}
