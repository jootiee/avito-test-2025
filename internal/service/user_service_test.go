package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *UserServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *UserServiceTestSuite) TestSetActive() {
	type args struct {
		userID   string
		isActive bool
	}
	type setupResult struct {
		svc      *UserService
		userRepo *mockUserRepository
	}

	tests := []struct {
		name       string
		args       args
		setup      func() setupResult
		assertFunc func(t *testing.T, user *domain.User, err error, sr setupResult)
	}{
		{
			name: "success",
			args: args{userID: "user1", isActive: false},
			setup: func() setupResult {
				userRepo := newMockUserRepository()
				_ = userRepo.UpsertUser(context.Background(), &domain.User{UserID: "user1", Username: "alice", TeamName: "backend", IsActive: true})
				return setupResult{svc: NewUserService(userRepo, newMockPullRequestRepository()), userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, user *domain.User, err error, sr setupResult) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.False(t, user.IsActive)
				// Verify in repository
				storedUser, _ := sr.userRepo.GetUser(context.Background(), "user1")
				assert.False(t, storedUser.IsActive)
			},
		},
		{
			name: "user not found",
			args: args{userID: "nonexistent", isActive: false},
			setup: func() setupResult {
				userRepo := newMockUserRepository()
				return setupResult{svc: NewUserService(userRepo, newMockPullRequestRepository()), userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, user *domain.User, err error, _ setupResult) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, "user not found", err.Error())
			},
		},
		{
			name: "repository error",
			args: args{userID: "user1", isActive: false},
			setup: func() setupResult {
				userRepo := newMockUserRepository()
				_ = userRepo.UpsertUser(context.Background(), &domain.User{UserID: "user1", Username: "alice", TeamName: "backend", IsActive: true})
				userRepo.setActiveErr = errors.New("database error")
				return setupResult{svc: NewUserService(userRepo, newMockPullRequestRepository()), userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, user *domain.User, err error, _ setupResult) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			sr := tt.setup()
			user, err := sr.svc.SetActive(s.ctx, tt.args.userID, tt.args.isActive)
			tt.assertFunc(s.T(), user, err, sr)
		})
	}
}

func (s *UserServiceTestSuite) TestGet() {
	type setupResult struct {
		svc      *UserService
		userRepo *mockUserRepository
	}

	tests := []struct {
		name       string
		userID     string
		setup      func() setupResult
		assertFunc func(t *testing.T, user *domain.User, err error)
	}{
		{
			name:   "success",
			userID: "user1",
			setup: func() setupResult {
				userRepo := newMockUserRepository()
				_ = userRepo.UpsertUser(context.Background(), &domain.User{UserID: "user1", Username: "alice", TeamName: "backend", IsActive: true})
				return setupResult{svc: NewUserService(userRepo, newMockPullRequestRepository()), userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, user *domain.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "user1", user.UserID)
				assert.Equal(t, "alice", user.Username)
			},
		},
		{
			name:   "not found",
			userID: "nonexistent",
			setup: func() setupResult {
				userRepo := newMockUserRepository()
				return setupResult{svc: NewUserService(userRepo, newMockPullRequestRepository()), userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, "user not found", err.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			sr := tt.setup()
			user, err := sr.svc.Get(s.ctx, tt.userID)
			tt.assertFunc(s.T(), user, err)
		})
	}
}

func (s *UserServiceTestSuite) TestGetReviews() {
	type setupResult struct {
		svc    *UserService
		prRepo *mockPullRequestRepository
	}

	tests := []struct {
		name       string
		userID     string
		setup      func() setupResult
		assertFunc func(t *testing.T, prs []*domain.PullRequest, err error)
	}{
		{
			name:   "success with multiple PRs",
			userID: "user1",
			setup: func() setupResult {
				prRepo := newMockPullRequestRepository()
				ctx := context.Background()
				_ = prRepo.CreatePR(ctx, &domain.PullRequest{PullRequestID: "pr1", PullRequestName: "Feature A", AuthorID: "author1", Status: domain.PullRequestStatusOpen, AssignedReviewers: []string{"user1", "user2"}})
				_ = prRepo.CreatePR(ctx, &domain.PullRequest{PullRequestID: "pr2", PullRequestName: "Feature B", AuthorID: "author2", Status: domain.PullRequestStatusOpen, AssignedReviewers: []string{"user1"}})
				_ = prRepo.CreatePR(ctx, &domain.PullRequest{PullRequestID: "pr3", PullRequestName: "Feature C", AuthorID: "author1", Status: domain.PullRequestStatusOpen, AssignedReviewers: []string{"user3"}})
				return setupResult{svc: NewUserService(newMockUserRepository(), prRepo), prRepo: prRepo}
			},
			assertFunc: func(t *testing.T, prs []*domain.PullRequest, err error) {
				assert.NoError(t, err)
				assert.Len(t, prs, 2)
				prIDs := make(map[string]bool)
				for _, pr := range prs {
					prIDs[pr.PullRequestID] = true
				}
				assert.True(t, prIDs["pr1"], "expected pr1 in results")
				assert.True(t, prIDs["pr2"], "expected pr2 in results")
			},
		},
		{
			name:   "empty result",
			userID: "user1",
			setup: func() setupResult {
				prRepo := newMockPullRequestRepository()
				return setupResult{svc: NewUserService(newMockUserRepository(), prRepo), prRepo: prRepo}
			},
			assertFunc: func(t *testing.T, prs []*domain.PullRequest, err error) {
				assert.NoError(t, err)
				assert.Empty(t, prs)
			},
		},
		{
			name:   "repository error",
			userID: "user1",
			setup: func() setupResult {
				prRepo := newMockPullRequestRepository()
				prRepo.getErr = errors.New("database error")
				return setupResult{svc: NewUserService(newMockUserRepository(), prRepo), prRepo: prRepo}
			},
			assertFunc: func(t *testing.T, prs []*domain.PullRequest, err error) {
				assert.Error(t, err)
				assert.Nil(t, prs)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			sr := tt.setup()
			prs, err := sr.svc.GetReviews(s.ctx, tt.userID)
			tt.assertFunc(s.T(), prs, err)
		})
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
