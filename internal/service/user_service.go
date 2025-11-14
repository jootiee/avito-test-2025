package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo UserRepository
	prRepo   PRRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo UserRepository, prRepo PRRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

// SetUserActive updates a user's active status
func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.userRepo.SetUserActive(ctx, userID, isActive); err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// GetUserReviews retrieves all PRs where user is assigned as reviewer
func (s *UserService) GetUserReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	prs, err := s.prRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
