package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/repository"
)

// Handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
	prRepo   repository.PRRepository
}

// Creates a new user service
func NewUserService(userRepo repository.UserRepository, prRepo repository.PRRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

// Updates a user's active status
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

// Retrieves a user by ID
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

// Retrieves all PRs where user is assigned as reviewer
func (s *UserService) GetUserReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	prs, err := s.prRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
