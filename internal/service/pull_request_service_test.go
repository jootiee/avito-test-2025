package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type PullRequestServiceTestSuite struct {
	suite.Suite
	ctx             context.Context
	pullRequestRepo *mockPullRequestRepository
	userRepo        *mockUserRepository
	teamRepo        *mockTeamRepository
	service         *PullRequestService
}

func (s *PullRequestServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.pullRequestRepo = newMockPullRequestRepository()
	s.userRepo = newMockUserRepository()
	s.teamRepo = newMockTeamRepository()
	s.service = NewPullRequestService(s.pullRequestRepo, s.userRepo, s.teamRepo)
}

func (s *PullRequestServiceTestSuite) setupTestData() {
	team := &domain.Team{
		TeamName: "backend",
		Members: []domain.User{
			{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
			{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
			{UserID: "user3", Username: "charlie", TeamName: "backend", IsActive: true},
			{UserID: "user4", Username: "dave", TeamName: "backend", IsActive: true},
			{UserID: "user5", Username: "eve", TeamName: "backend", IsActive: false},
		},
	}
	s.teamRepo.CreateTeam(s.ctx, team)

	for _, user := range team.Members {
		userCopy := user
		s.userRepo.UpsertUser(s.ctx, &userCopy)
	}
}

func (s *PullRequestServiceTestSuite) TestCreate() {
	testCases := []struct {
		name          string
		prID          string
		prName        string
		authorID      string
		setupData     func()
		expectedError string
		validate      func(pr *domain.PullRequest)
	}{
		{
			name:      "success",
			prID:      "pr1",
			prName:    "Add feature",
			authorID:  "author1",
			setupData: s.setupTestData,
			validate: func(pr *domain.PullRequest) {
				assert.Equal(s.T(), "pr1", pr.PullRequestID)
				assert.Equal(s.T(), domain.PullRequestStatusOpen, pr.Status)
				assert.Equal(s.T(), "author1", pr.AuthorID)
				assert.LessOrEqual(s.T(), len(pr.AssignedReviewers), 2)
				for _, reviewer := range pr.AssignedReviewers {
					assert.NotEqual(s.T(), "author1", reviewer, "author should not be assigned as reviewer")
				}
			},
		},
		{
			name:     "PR already exists",
			prID:     "pr1",
			prName:   "Add feature",
			authorID: "author1",
			setupData: func() {
				s.setupTestData()
				s.service.Create(s.ctx, "pr1", "Add feature", "author1")
			},
			expectedError: "PR already exists",
		},
		{
			name:          "author not found",
			prID:          "pr1",
			prName:        "Add feature",
			authorID:      "nonexistent",
			setupData:     func() {},
			expectedError: "user not found",
		},
		{
			name:     "team not found",
			prID:     "pr1",
			prName:   "Add feature",
			authorID: "author1",
			setupData: func() {
				user := &domain.User{
					UserID:   "author1",
					Username: "alice",
					TeamName: "nonexistent",
					IsActive: true,
				}
				s.userRepo.UpsertUser(s.ctx, user)
			},
			expectedError: "team not found",
		},
		{
			name:     "no active members to assign",
			prID:     "pr1",
			prName:   "Add feature",
			authorID: "author1",
			setupData: func() {
				team := &domain.Team{
					TeamName: "backend",
					Members: []domain.User{
						{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
						{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: false},
					},
				}
				s.teamRepo.CreateTeam(s.ctx, team)
				for _, user := range team.Members {
					userCopy := user
					s.userRepo.UpsertUser(s.ctx, &userCopy)
				}
			},
			validate: func(pr *domain.PullRequest) {
				assert.Equal(s.T(), "pr1", pr.PullRequestID)
				assert.Empty(s.T(), pr.AssignedReviewers)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			if tc.setupData != nil {
				tc.setupData()
			}

			pr, err := s.service.Create(s.ctx, tc.prID, tc.prName, tc.authorID)

			if tc.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tc.expectedError)
			} else {
				s.NoError(err)
				if tc.validate != nil {
					tc.validate(pr)
				}
			}
		})
	}
}

func (s *PullRequestServiceTestSuite) TestMerge() {
	testCases := []struct {
		name          string
		prID          string
		setupData     func()
		expectedError string
		validate      func(pr *domain.PullRequest)
	}{
		{
			name: "success",
			prID: "pr1",
			setupData: func() {
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			validate: func(pr *domain.PullRequest) {
				assert.Equal(s.T(), domain.PullRequestStatusMerged, pr.Status)
				assert.NotNil(s.T(), pr.MergedAt)
			},
		},
		{
			name: "idempotent merge",
			prID: "pr1",
			setupData: func() {
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				pr.Status = domain.PullRequestStatusMerged
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			validate: func(pr *domain.PullRequest) {
				assert.Equal(s.T(), domain.PullRequestStatusMerged, pr.Status)
			},
		},
		{
			name:          "PR not found",
			prID:          "nonexistent",
			setupData:     func() {},
			expectedError: "PR not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			if tc.setupData != nil {
				tc.setupData()
			}

			pr, err := s.service.Merge(s.ctx, tc.prID)

			if tc.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tc.expectedError)
			} else {
				s.NoError(err)
				if tc.validate != nil {
					tc.validate(pr)
				}
			}
		})
	}
}

func (s *PullRequestServiceTestSuite) TestReassignReviewer() {
	testCases := []struct {
		name          string
		prID          string
		oldUserID     string
		setupData     func()
		expectedError string
		validate      func(pr *domain.PullRequest, newUserID string)
	}{
		{
			name:      "success",
			prID:      "pr1",
			oldUserID: "user2",
			setupData: func() {
				s.setupTestData()
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				pr.AssignedReviewers = []string{"user2", "user3"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			validate: func(pr *domain.PullRequest, newUserID string) {
				assert.NotContains(s.T(), pr.AssignedReviewers, "user2")
				assert.Contains(s.T(), pr.AssignedReviewers, newUserID)
				assert.Len(s.T(), pr.AssignedReviewers, 2)
			},
		},
		{
			name:          "PR not found",
			prID:          "nonexistent",
			oldUserID:     "user2",
			setupData:     func() { s.setupTestData() },
			expectedError: "PR not found",
		},
		{
			name:      "PR already merged",
			prID:      "pr1",
			oldUserID: "user2",
			setupData: func() {
				s.setupTestData()
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				pr.Status = domain.PullRequestStatusMerged
				pr.AssignedReviewers = []string{"user2", "user3"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			expectedError: "cannot reassign on merged PullRequest",
		},
		{
			name:      "user not assigned",
			prID:      "pr1",
			oldUserID: "user4",
			setupData: func() {
				s.setupTestData()
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				pr.AssignedReviewers = []string{"user2", "user3"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			expectedError: "reviewer is not assigned to this PR",
		},
		{
			name:      "old reviewer not found",
			prID:      "pr1",
			oldUserID: "nonexistent",
			setupData: func() {
				s.setupTestData()
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				pr.AssignedReviewers = []string{"user2"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			expectedError: "reviewer is not assigned to this PR",
		},
		{
			name:      "no available replacement",
			prID:      "pr1",
			oldUserID: "user2",
			setupData: func() {
				team := &domain.Team{
					TeamName: "backend",
					Members: []domain.User{
						{UserID: "author1", Username: "alice", TeamName: "backend", IsActive: true},
						{UserID: "user2", Username: "bob", TeamName: "backend", IsActive: true},
					},
				}
				_ = s.teamRepo.CreateTeam(s.ctx, team)
				for _, user := range team.Members {
					userCopy := user
					_ = s.userRepo.UpsertUser(s.ctx, &userCopy)
				}
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				pr.AssignedReviewers = []string{"user2"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
			expectedError: "no active replacement candidate in team",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			if tc.setupData != nil {
				tc.setupData()
			}

			pr, newUserID, err := s.service.ReassignReviewer(s.ctx, tc.prID, tc.oldUserID)

			if tc.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tc.expectedError)
				return
			}

			s.NoError(err)
			if tc.validate != nil {
				tc.validate(pr, newUserID)
			}
		})
	}
}

func (s *PullRequestServiceTestSuite) TestGet() {
	testCases := []struct {
		name          string
		prID          string
		setupData     func()
		expectedError string
	}{
		{
			name: "success",
			prID: "pr1",
			setupData: func() {
				pr := domain.NewPullRequest("pr1", "Add feature", "author1")
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr)
			},
		},
		{
			name:          "not found",
			prID:          "nonexistent",
			setupData:     func() {},
			expectedError: "PR not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			if tc.setupData != nil {
				tc.setupData()
			}

			pr, err := s.service.Get(s.ctx, tc.prID)

			if tc.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tc.expectedError)
			} else {
				s.NoError(err)
				assert.Equal(s.T(), tc.prID, pr.PullRequestID)
			}
		})
	}
}

func (s *PullRequestServiceTestSuite) TestGetStatistics() {
	testCases := []struct {
		name      string
		setupData func()
		repoErr   error
		validate  func(userCounts, prCounts map[string]int)
	}{
		{
			name: "with multiple PRs",
			setupData: func() {
				pr1 := domain.NewPullRequest("pr1", "Feature A", "author1")
				pr1.AssignedReviewers = []string{"user2", "user3"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr1)

				pr2 := domain.NewPullRequest("pr2", "Feature B", "author2")
				pr2.AssignedReviewers = []string{"user2", "user4"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr2)

				pr3 := domain.NewPullRequest("pr3", "Feature C", "author3")
				pr3.AssignedReviewers = []string{"user3"}
				_ = s.pullRequestRepo.CreatePR(s.ctx, pr3)
			},
			validate: func(userCounts, prCounts map[string]int) {
				assert.Equal(s.T(), 2, userCounts["user2"])
				assert.Equal(s.T(), 2, userCounts["user3"])
				assert.Equal(s.T(), 1, userCounts["user4"])
				assert.Equal(s.T(), 2, prCounts["pr1"])
				assert.Equal(s.T(), 2, prCounts["pr2"])
				assert.Equal(s.T(), 1, prCounts["pr3"])
			},
		},
		{
			name:      "empty data",
			setupData: func() {},
			validate: func(userCounts, prCounts map[string]int) {
				assert.Empty(s.T(), userCounts)
				assert.Empty(s.T(), prCounts)
			},
		},
		{
			name:      "repository error",
			setupData: func() {},
			repoErr:   errors.New("database error"),
			validate:  nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			if tc.setupData != nil {
				tc.setupData()
			}
			if tc.repoErr != nil {
				s.pullRequestRepo.getErr = tc.repoErr
			}

			userCounts, prCounts, err := s.service.GetStatistics(s.ctx)

			if tc.repoErr != nil {
				s.Error(err)
				assert.Equal(s.T(), tc.repoErr.Error(), err.Error())
				return
			}

			s.NoError(err)
			if tc.validate != nil {
				tc.validate(userCounts, prCounts)
			}
		})
	}
}

func TestPullRequestService(t *testing.T) {
	suite.Run(t, new(PullRequestServiceTestSuite))
}
