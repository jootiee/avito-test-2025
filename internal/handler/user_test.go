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

type UserHandlerTestSuite struct {
	suite.Suite
}

func (s *UserHandlerTestSuite) TestUserEndpoints() {
	type args struct {
		method string
		path   string
		body   interface{}
	}
	type setupResult struct {
		h        *Handler
		userRepo *mockUserRepo
		prRepo   *mockPRRepo
	}

	tests := []struct {
		name       string
		args       args
		setup      func() setupResult
		before     func(sr setupResult)
		assertFunc func(w *httptest.ResponseRecorder, sr setupResult)
	}{
		{
			name: "set user active success",
			args: args{method: http.MethodPost, path: "/users/setIsActive", body: dto.SetActiveRequest{UserID: "user1", IsActive: false}},
			setup: func() setupResult {
				userRepo := &mockUserRepo{users: map[string]*domain.User{"user1": {UserID: "user1", Username: "alice", TeamName: "backend", IsActive: true}}}
				sr := setupResult{userRepo: userRepo, prRepo: &mockPRRepo{prs: map[string]*domain.PullRequest{}}}
				sr.h = New(service.New(&mockTeamRepo{teams: map[string]*domain.Team{}}, userRepo, sr.prRepo), &mockLogger{})
				return sr
			},
			assertFunc: func(w *httptest.ResponseRecorder, _ setupResult) {
				assert.Equal(s.T(), http.StatusOK, w.Code)
				var resp dto.UserResponse
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.False(s.T(), resp.User.IsActive)
			},
		},
		{
			name: "set user active not found",
			args: args{method: http.MethodPost, path: "/users/setIsActive", body: dto.SetActiveRequest{UserID: "missing", IsActive: false}},
			setup: func() setupResult {
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				sr := setupResult{userRepo: userRepo, prRepo: &mockPRRepo{prs: map[string]*domain.PullRequest{}}}
				sr.h = New(service.New(&mockTeamRepo{teams: map[string]*domain.Team{}}, userRepo, sr.prRepo), &mockLogger{})
				return sr
			},
			assertFunc: func(w *httptest.ResponseRecorder, _ setupResult) {
				assert.Equal(s.T(), http.StatusNotFound, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "user not found", resp.Error.Message)
			},
		},
		{
			name: "set user active invalid JSON",
			args: args{method: http.MethodPost, path: "/users/setIsActive", body: []byte("invalid")},
			setup: func() setupResult {
				userRepo := &mockUserRepo{users: map[string]*domain.User{}}
				sr := setupResult{userRepo: userRepo, prRepo: &mockPRRepo{prs: map[string]*domain.PullRequest{}}}
				sr.h = New(service.New(&mockTeamRepo{teams: map[string]*domain.Team{}}, userRepo, sr.prRepo), &mockLogger{})
				return sr
			},
			assertFunc: func(w *httptest.ResponseRecorder, _ setupResult) {
				assert.Equal(s.T(), http.StatusBadRequest, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), dto.ErrCodeNotFound, resp.Error.Code)
			},
		},
		{
			name: "get user reviews success",
			args: args{method: http.MethodGet, path: "/users/getReview?user_id=user1"},
			setup: func() setupResult {
				userRepo := &mockUserRepo{users: map[string]*domain.User{"user1": {UserID: "user1", Username: "alice", TeamName: "backend", IsActive: true}}}
				prRepo := &mockPRRepo{prs: map[string]*domain.PullRequest{}}
				sr := setupResult{userRepo: userRepo, prRepo: prRepo}
				sr.h = New(service.New(&mockTeamRepo{teams: map[string]*domain.Team{}}, userRepo, prRepo), &mockLogger{})
				return sr
			},
			before: func(sr setupResult) {
				pr := domain.NewPullRequest("pr1", "Feature", "author1")
				pr.AssignedReviewers = []string{"user1"}
				_ = sr.prRepo.CreatePR(context.Background(), pr)
			},
			assertFunc: func(w *httptest.ResponseRecorder, _ setupResult) {
				assert.Equal(s.T(), http.StatusOK, w.Code)
				var resp dto.UserReviewsResponse
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "user1", resp.UserID)
				assert.Len(s.T(), resp.PullRequests, 1)
				assert.Equal(s.T(), "pr1", resp.PullRequests[0].PullRequestID)
			},
		},
		{
			name: "get user reviews missing user_id",
			args: args{method: http.MethodGet, path: "/users/getReview"},
			setup: func() setupResult {
				sr := setupResult{userRepo: &mockUserRepo{users: map[string]*domain.User{}}, prRepo: &mockPRRepo{prs: map[string]*domain.PullRequest{}}}
				sr.h = New(service.New(&mockTeamRepo{teams: map[string]*domain.Team{}}, sr.userRepo, sr.prRepo), &mockLogger{})
				return sr
			},
			assertFunc: func(w *httptest.ResponseRecorder, _ setupResult) {
				assert.Equal(s.T(), http.StatusBadRequest, w.Code)
				var resp dto.APIError
				assert.NoError(s.T(), json.NewDecoder(w.Body).Decode(&resp))
				assert.Equal(s.T(), "user_id required", resp.Error.Message)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			sr := tc.setup()
			if tc.before != nil {
				tc.before(sr)
			}

			var bodyBytes []byte
			switch b := tc.args.body.(type) {
			case nil:
			case []byte:
				bodyBytes = b
			default:
				marshaled, err := json.Marshal(b)
				s.NoError(err)
				bodyBytes = marshaled
			}
			req := httptest.NewRequest(tc.args.method, tc.args.path, bytes.NewReader(bodyBytes))
			if bodyBytes != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			sr.h.ServeHTTP(w, req)
			tc.assertFunc(w, sr)
		})
	}
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
