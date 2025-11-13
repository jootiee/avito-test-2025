package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/repository"
)

// Handles pull request business logic
type PRService struct {
	prRepo   repository.PRRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

// Creates a new PR service
func NewPRService(prRepo repository.PRRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// Creates a new pull request and auto-assigns reviewers
func (s *PRService) CreatePR(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	exists, err := s.prRepo.PRExists(ctx, prID)
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

	if err := s.prRepo.CreatePR(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// Marks a PR as merged
func (s *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("PR not found")
		}
		return nil, err
	}

	// If already merged, return as-is
	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	// Update status
	now := time.Now().UTC()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	if err := s.prRepo.UpdatePR(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// Replaces a reviewer with another from the same team
func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("PR not found")
		}
		return nil, "", err
	}

	// Check if PR is merged
	if pr.Status == domain.PRStatusMerged {
		return nil, "", errors.New("cannot reassign on merged PR")
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
	if err := s.prRepo.UpdatePR(ctx, pr); err != nil {
		return nil, "", err
	}

	return pr, newUserID, nil
}

// GetPR retrieves a PR by ID
func (s *PRService) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("PR not found")
		}
		return nil, err
	}
	return pr, nil
}
