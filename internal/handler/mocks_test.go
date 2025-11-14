package handler

import (
	"context"
	"errors"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

// Mock services for handler testing

// mockTeamService is a mock implementation of TeamService
type mockTeamService struct {
	createTeamFunc func(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error)
	getTeamFunc    func(ctx context.Context, teamName string) (*domain.Team, error)
}

func (m *mockTeamService) CreateTeam(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error) {
	if m.createTeamFunc != nil {
		return m.createTeamFunc(ctx, teamName, members)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	if m.getTeamFunc != nil {
		return m.getTeamFunc(ctx, teamName)
	}
	return nil, errors.New("not implemented")
}

// mockUserService is a mock implementation of UserService
type mockUserService struct {
	setUserActiveFunc  func(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	getUserFunc        func(ctx context.Context, userID string) (*domain.User, error)
	getUserReviewsFunc func(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}

func (m *mockUserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if m.setUserActiveFunc != nil {
		return m.setUserActiveFunc(ctx, userID, isActive)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetUserReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	if m.getUserReviewsFunc != nil {
		return m.getUserReviewsFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

// mockPRService is a mock implementation of PRService
type mockPRService struct {
	createPRFunc         func(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error)
	mergePRFunc          func(ctx context.Context, prID string) (*domain.PullRequest, error)
	reassignReviewerFunc func(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error)
	getPRFunc            func(ctx context.Context, prID string) (*domain.PullRequest, error)
}

func (m *mockPRService) CreatePR(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	if m.createPRFunc != nil {
		return m.createPRFunc(ctx, prID, prName, authorID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockPRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if m.mergePRFunc != nil {
		return m.mergePRFunc(ctx, prID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockPRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	if m.reassignReviewerFunc != nil {
		return m.reassignReviewerFunc(ctx, prID, oldUserID)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockPRService) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if m.getPRFunc != nil {
		return m.getPRFunc(ctx, prID)
	}
	return nil, errors.New("not implemented")
}

// mockLogger is a mock implementation of logger.Interface
type mockLogger struct{}

func (m *mockLogger) Debug(msg interface{}, keysAndValues ...interface{}) {}
func (m *mockLogger) Info(msg string, keysAndValues ...interface{})       {}
func (m *mockLogger) Warn(msg string, keysAndValues ...interface{})       {}
func (m *mockLogger) Error(msg interface{}, keysAndValues ...interface{}) {}
func (m *mockLogger) Fatal(msg interface{}, keysAndValues ...interface{}) {}

// testServiceBuilder helps build service with mocked behaviors
type testServiceBuilder struct {
	teamService *mockTeamService
	userService *mockUserService
	prService   *mockPRService
}

func newTestService() *testServiceBuilder {
	return &testServiceBuilder{
		teamService: &mockTeamService{},
		userService: &mockUserService{},
		prService:   &mockPRService{},
	}
}

func (b *testServiceBuilder) withTeamCreate(fn func(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error)) *testServiceBuilder {
	b.teamService.createTeamFunc = fn
	return b
}

func (b *testServiceBuilder) withTeamGet(fn func(ctx context.Context, teamName string) (*domain.Team, error)) *testServiceBuilder {
	b.teamService.getTeamFunc = fn
	return b
}

func (b *testServiceBuilder) withUserSetActive(fn func(ctx context.Context, userID string, isActive bool) (*domain.User, error)) *testServiceBuilder {
	b.userService.setUserActiveFunc = fn
	return b
}

func (b *testServiceBuilder) withUserReviews(fn func(ctx context.Context, userID string) ([]*domain.PullRequest, error)) *testServiceBuilder {
	b.userService.getUserReviewsFunc = fn
	return b
}

func (b *testServiceBuilder) withPRCreate(fn func(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error)) *testServiceBuilder {
	b.prService.createPRFunc = fn
	return b
}

func (b *testServiceBuilder) withPRMerge(fn func(ctx context.Context, prID string) (*domain.PullRequest, error)) *testServiceBuilder {
	b.prService.mergePRFunc = fn
	return b
}

func (b *testServiceBuilder) withPRReassign(fn func(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error)) *testServiceBuilder {
	b.prService.reassignReviewerFunc = fn
	return b
}
