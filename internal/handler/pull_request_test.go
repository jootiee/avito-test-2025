package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
)

type PullRequestHandlerTestSuite struct {
	suite.Suite
}

func (s *PullRequestHandlerTestSuite) mkTeam() *domain.Team {
	return &domain.Team{TeamName: "backend", Members: []domain.User{
		{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
		{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
		{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
	}}
}

func (s *PullRequestHandlerTestSuite) TestEndpoints() {
	type testCase struct {
		name       string
		method     string
		path       string
		body       interface{}
		setup      func() *Handler
		assertFunc func(w *httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			name:   "create PullRequest success",
			method: http.MethodPost,
			path:   "/pullRequest/create",
			body:   dto.CreatePullRequestRequest{PullRequestID: "pr1", PullRequestName: "Add feature", AuthorID: "author1"},
			setup: func() *Handler {
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": s.mkTeam()}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				for _, m := range s.mkTeam().Members {
					userRepo.UpsertUser(context.Background(), &m)
				}
				h := New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
				return h
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusCreated, w.Code)
				var resp dto.PullRequestResponse
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "pr1", resp.PullRequest.PullRequestID)
				assert.Equal(s.T(), domain.PullRequestStatusOpen, resp.PullRequest.Status)
			},
		},
		{
			name:   "create PR duplicate",
			method: http.MethodPost,
			path:   "/pullRequest/create",
			body:   dto.CreatePullRequestRequest{PullRequestID: "pr1", PullRequestName: "Add feature", AuthorID: "author1"},
			setup: func() *Handler {
				team := s.mkTeam()
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": team}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				for _, m := range team.Members {
					userRepo.UpsertUser(context.Background(), &m)
				}
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				prRepo := &mockPRRepo{prs: map[string]*domain.PullRequest{"pr1": pr}}
				return New(service.New(teamRepo, userRepo, prRepo), &mockLogger{})
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusConflict, w.Code)
			},
		},
		{
			name:   "merge PR success",
			method: http.MethodPost,
			path:   "/pullRequest/merge",
			body:   dto.MergeRequest{PullRequestID: "pr1"},
			setup: func() *Handler {
				team := s.mkTeam()
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": team}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				for _, m := range team.Members {
					userRepo.UpsertUser(context.Background(), &m)
				}
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				prRepo := &mockPRRepo{prs: map[string]*domain.PullRequest{"pr1": pr}}
				return New(service.New(teamRepo, userRepo, prRepo), &mockLogger{})
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusOK, w.Code)
				var resp dto.PullRequestResponse
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), domain.PullRequestStatusMerged, resp.PullRequest.Status)
			},
		},
		{
			name:   "reassign reviewer merged PR conflict",
			method: http.MethodPost,
			path:   "/pullRequest/reassign",
			body:   dto.ReassignRequest{PullRequestID: "pr1", OldUserID: "user2"},
			setup: func() *Handler {
				team := s.mkTeam()
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": team}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				for _, m := range team.Members {
					userRepo.UpsertUser(context.Background(), &m)
				}
				pr := domain.NewPullRequest("pr1", "Add feature", "user1")
				pr.AssignedReviewers = []string{"user2"}
				pr.Status = domain.PullRequestStatusMerged
				prRepo := &mockPRRepo{prs: map[string]*domain.PullRequest{"pr1": pr}}
				return New(service.New(teamRepo, userRepo, prRepo), &mockLogger{})
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusConflict, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), dto.ErrCodePullRequestMerged, resp.Error.Code)
			},
		},
		{
			name:   "create PR author not found",
			method: http.MethodPost,
			path:   "/pullRequest/create",
			body:   dto.CreatePullRequestRequest{PullRequestID: "pr2", PullRequestName: "Another feature", AuthorID: "nonexistent"},
			setup: func() *Handler {
				team := s.mkTeam()
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": team}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				for _, m := range team.Members {
					userRepo.UpsertUser(context.Background(), &m)
				}
				return New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusNotFound, w.Code)
			},
		},
		{
			name:   "merge PR not found",
			method: http.MethodPost,
			path:   "/pullRequest/merge",
			body:   dto.MergeRequest{PullRequestID: "nonexistent"},
			setup: func() *Handler {
				team := s.mkTeam()
				teamRepo := &mockTeamRepo{teams: map[string]*domain.Team{"backend": team}}
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				for _, m := range team.Members {
					userRepo.UpsertUser(context.Background(), &m)
				}
				return New(service.New(teamRepo, userRepo, &mockPRRepo{prs: map[string]*domain.PullRequest{}}), &mockLogger{})
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusNotFound, w.Code)
			},
		},
		{
			name:   "create PR invalid JSON",
			method: http.MethodPost,
			path:   "/pullRequest/create",
			body:   []byte("invalid"),
			setup: func() *Handler {
				return New(&service.Service{}, &mockLogger{})
			},
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusBadRequest, w.Code)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			h := tc.setup()
			var reqBody []byte
			switch b := tc.body.(type) {
			case nil:
			case []byte:
				reqBody = b
			default:
				marshalled, err := json.Marshal(b)
				s.NoError(err)
				reqBody = marshalled
			}
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(reqBody))
			if reqBody != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			tc.assertFunc(w)
		})
	}
}

func TestPRHandler(t *testing.T) {
	suite.Run(t, new(PullRequestHandlerTestSuite))
}
