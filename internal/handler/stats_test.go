package handler

import (
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

type StatsHandlerTestSuite struct {
	suite.Suite
}

func (s *StatsHandlerTestSuite) TestStatsEndpoints() {
	type setupResult struct {
		h http.Handler
	}

	testCases := []struct {
		name             string
		setup            func() setupResult
		expectedCode     int
		validateResponse func(resp dto.StatsResponse)
	}{
		{
			name: "valid with multiple PRs",
			setup: func() setupResult {
				teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
				userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
				prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

				pr1 := domain.NewPullRequest("pr1", "Feature A", "author1")
				pr1.AssignedReviewers = []string{"user2", "user3"}
				_ = prRepo.CreatePR(context.Background(), pr1)

				pr2 := domain.NewPullRequest("pr2", "Feature B", "author2")
				pr2.AssignedReviewers = []string{"user2", "user4"}
				_ = prRepo.CreatePR(context.Background(), pr2)

				pr3 := domain.NewPullRequest("pr3", "Feature C", "author3")
				pr3.AssignedReviewers = []string{"user3"}
				_ = prRepo.CreatePR(context.Background(), pr3)

				svc := service.New(teamRepo, userRepo, prRepo)
				h := New(svc, &mockLogger{})
				return setupResult{h: h}
			},
			expectedCode: http.StatusOK,
			validateResponse: func(resp dto.StatsResponse) {
				assert.Len(s.T(), resp.UserStats, 3)
				assert.Len(s.T(), resp.PRStats, 3)

				userCountMap := make(map[string]int)
				for _, userStat := range resp.UserStats {
					userCountMap[userStat.UserID] = userStat.AssignmentCount
				}

				assert.Equal(s.T(), 2, userCountMap["user2"])
				assert.Equal(s.T(), 2, userCountMap["user3"])
				assert.Equal(s.T(), 1, userCountMap["user4"])

				prCountMap := make(map[string]int)
				for _, prStat := range resp.PRStats {
					prCountMap[prStat.PullRequestID] = prStat.ReviewerCount
				}

				assert.Equal(s.T(), 2, prCountMap["pr1"])
				assert.Equal(s.T(), 2, prCountMap["pr2"])
				assert.Equal(s.T(), 1, prCountMap["pr3"])
			},
		},
		{
			name: "empty data",
			setup: func() setupResult {
				teamRepo := &mockTeamRepo{teams: make(map[string]*domain.Team)}
				userRepo := &mockUserRepo{users: make(map[string]*domain.User)}
				prRepo := &mockPRRepo{prs: make(map[string]*domain.PullRequest)}

				svc := service.New(teamRepo, userRepo, prRepo)
				h := New(svc, &mockLogger{})
				return setupResult{h: h}
			},
			expectedCode: http.StatusOK,
			validateResponse: func(resp dto.StatsResponse) {
				assert.Empty(s.T(), resp.UserStats)
				assert.Empty(s.T(), resp.PRStats)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sr := tc.setup()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/stats", http.NoBody)
			sr.h.ServeHTTP(rec, req)

			assert.Equal(s.T(), tc.expectedCode, rec.Code)

			if tc.validateResponse != nil {
				var resp dto.StatsResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				tc.validateResponse(resp)
			}
		})
	}
}

func TestStatsHandler(t *testing.T) {
	suite.Run(t, new(StatsHandlerTestSuite))
}
