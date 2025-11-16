package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

// PullRequestRepository defines the interface for pull request data access
type PullRequestRepository interface {
	CreatePR(ctx context.Context, pr *domain.PullRequest) error
	GetPR(ctx context.Context, prID string) (*domain.PullRequest, error)
	UpdatePR(ctx context.Context, pr *domain.PullRequest) error
	PRExists(ctx context.Context, prID string) (bool, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error)
	GetUserAssignmentCounts(ctx context.Context) (map[string]int, error)
	GetPRReviewerCounts(ctx context.Context) (map[string]int, error)
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	UpsertUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)
}

// TeamRepository defines the interface for team data access
type TeamRepository interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}

// PullRequestService handles pull request business logic
type PullRequestService struct {
	pullRequestRepo PullRequestRepository
	userRepo        UserRepository
	teamRepo        TeamRepository
}

// NewPullRequestService creates a new PullRequest service
func NewPullRequestService(
	pullRequestRepo PullRequestRepository,
	userRepo UserRepository,
	teamRepo TeamRepository,
) *PullRequestService {
	return &PullRequestService{
		pullRequestRepo: pullRequestRepo,
		userRepo:        userRepo,
		teamRepo:        teamRepo,
	}
}

// Create creates a new pull request and auto-assigns reviewers
func (s *PullRequestService) Create(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	exists, err := s.pullRequestRepo.PRExists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("PR already exists")
	}

	author, err := s.userRepo.GetUser(ctx, authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	team, err := s.teamRepo.GetTeam(ctx, author.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}

	var candidates []string
	for _, member := range team.Members {
		if member.UserID == authorID {
			continue
		}
		if member.IsActive {
			candidates = append(candidates, member.UserID)
		}
	}

	assigned := SelectRandomReviewers(candidates, 2)

	pr := domain.NewPullRequest(prID, prName, authorID)
	now := time.Now().UTC()
	pr.CreatedAt = &now
	pr.AssignedReviewers = assigned

	if err := s.pullRequestRepo.CreatePR(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// Merge marks a PR as merged
func (s *PullRequestService) Merge(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.pullRequestRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("PR not found")
		}
		return nil, err
	}

	// If already merged, return as-is
	if pr.Status == domain.PullRequestStatusMerged {
		return pr, nil
	}

	// Update status
	now := time.Now().UTC()
	pr.Status = domain.PullRequestStatusMerged
	pr.MergedAt = &now

	if err := s.pullRequestRepo.UpdatePR(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// ReassignReviewer replaces a reviewer with another from the same team
func (s *PullRequestService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	pr, err := s.pullRequestRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("PR not found")
		}
		return nil, "", err
	}

	// Check if PR is merged
	if pr.Status == domain.PullRequestStatusMerged {
		return nil, "", errors.New("cannot reassign on merged PullRequest")
	}

	// Check if old user is assigned
	found := false
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", errors.New("reviewer is not assigned to this PR")
	}

	// Get old user's team
	oldUser, err := s.userRepo.GetUser(ctx, oldUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	team, err := s.teamRepo.GetTeam(ctx, oldUser.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("team not found")
		}
		return nil, "", err
	}

	// Build set of already assigned reviewers
	assignedSet := make(map[string]struct{})
	for _, reviewerID := range pr.AssignedReviewers {
		assignedSet[reviewerID] = struct{}{}
	}

	// Collect candidates: active, not already assigned, not the old user, not the author
	var candidates []string
	for _, member := range team.Members {
		if !member.IsActive {
			continue
		}
		if member.UserID == oldUserID {
			continue
		}
		if member.UserID == pr.AuthorID {
			continue
		}
		if _, already := assignedSet[member.UserID]; already {
			continue
		}
		candidates = append(candidates, member.UserID)
	}

	if len(candidates) == 0 {
		return nil, "", errors.New("no active replacement candidate in team")
	}

	// Pick random candidate
	newReviewers := SelectRandomReviewers(candidates, 1)
	newUserID := newReviewers[0]

	// Replace old user with new user
	for i, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			pr.AssignedReviewers[i] = newUserID
			break
		}
	}

	// Update PR
	if err := s.pullRequestRepo.UpdatePR(ctx, pr); err != nil {
		return nil, "", err
	}

	return pr, newUserID, nil
}

// Get retrieves a PR by ID
func (s *PullRequestService) Get(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.pullRequestRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("PR not found")
		}
		return nil, err
	}
	return pr, nil
}

// GetStatistics retrieves reviewer assignment statistics
func (s *PullRequestService) GetStatistics(ctx context.Context) (userCounts, prCounts map[string]int, err error) {
	userCounts, err = s.pullRequestRepo.GetUserAssignmentCounts(ctx)
	if err != nil {
		return nil, nil, err
	}

	prCounts, err = s.pullRequestRepo.GetPRReviewerCounts(ctx)
	if err != nil {
		return nil, nil, err
	}

	return userCounts, prCounts, nil
}
