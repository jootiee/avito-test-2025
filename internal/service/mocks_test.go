package service

import (
	"context"
	"errors"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

// Mock repositories for testing

// mockTeamRepository is a mock implementation of TeamRepository
type mockTeamRepository struct {
	teams     map[string]*domain.Team
	createErr error
	getErr    error
	existsErr error
}

func newMockTeamRepository() *mockTeamRepository {
	return &mockTeamRepository{
		teams: make(map[string]*domain.Team),
	}
}

func (m *mockTeamRepository) CreateTeam(_ context.Context, team *domain.Team) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.teams[team.TeamName] = team
	return nil
}

func (m *mockTeamRepository) GetTeam(_ context.Context, teamName string) (*domain.Team, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	team, exists := m.teams[teamName]
	if !exists {
		return nil, errors.New("team not found")
	}
	return team, nil
}

func (m *mockTeamRepository) TeamExists(_ context.Context, teamName string) (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}
	_, exists := m.teams[teamName]
	return exists, nil
}

// mockUserRepository is a mock implementation of UserRepository
type mockUserRepository struct {
	users        map[string]*domain.User
	upsertErr    error
	getErr       error
	setActiveErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) UpsertUser(ctx context.Context, user *domain.User) error {
	if m.upsertErr != nil {
		return m.upsertErr
	}
	m.users[user.UserID] = user
	return nil
}

func (m *mockUserRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	user, exists := m.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) SetUserActive(ctx context.Context, userID string, isActive bool) error {
	if m.setActiveErr != nil {
		return m.setActiveErr
	}
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.IsActive = isActive
	return nil
}

func (m *mockUserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var members []domain.User
	for _, user := range m.users {
		if user.TeamName == teamName {
			members = append(members, *user)
		}
	}
	return members, nil
}

// mockPullRequestRepository is a mock implementation of PRRepository
type mockPullRequestRepository struct {
	prs       map[string]*domain.PullRequest
	createErr error
	getErr    error
	updateErr error
	existsErr error
}

func newMockPullRequestRepository() *mockPullRequestRepository {
	return &mockPullRequestRepository{
		prs: make(map[string]*domain.PullRequest),
	}
}

func (m *mockPullRequestRepository) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.prs[pr.PullRequestID] = pr
	return nil
}

func (m *mockPullRequestRepository) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	pr, exists := m.prs[prID]
	if !exists {
		return nil, errors.New("PR not found")
	}
	return pr, nil
}

func (m *mockPullRequestRepository) UpdatePR(ctx context.Context, pr *domain.PullRequest) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.prs[pr.PullRequestID]; !exists {
		return errors.New("PR not found")
	}
	m.prs[pr.PullRequestID] = pr
	return nil
}

func (m *mockPullRequestRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}
	_, exists := m.prs[prID]
	return exists, nil
}

func (m *mockPullRequestRepository) GetPRsByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
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

func (m *mockPullRequestRepository) GetUserAssignmentCounts(ctx context.Context) (map[string]int, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	counts := make(map[string]int)
	for _, pr := range m.prs {
		for _, reviewer := range pr.AssignedReviewers {
			counts[reviewer]++
		}
	}
	return counts, nil
}

func (m *mockPullRequestRepository) GetPRReviewerCounts(ctx context.Context) (map[string]int, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	counts := make(map[string]int)
	for prID, pr := range m.prs {
		counts[prID] = len(pr.AssignedReviewers)
	}
	return counts, nil
}
