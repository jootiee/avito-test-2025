package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

func TestTeamService_CreateTeam_Success(t *testing.T) {
	teamRepo := newMockTeamRepository()
	userRepo := newMockUserRepository()
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()
	teamName := "backend"
	members := []domain.User{
		{UserID: "user1", Username: "alice", IsActive: true},
		{UserID: "user2", Username: "bob", IsActive: true},
	}

	team, err := service.CreateTeam(ctx, teamName, members)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if team == nil {
		t.Fatal("expected team to be returned")
	}

	if team.TeamName != teamName {
		t.Errorf("expected team name %s, got %s", teamName, team.TeamName)
	}

	if len(team.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(team.Members))
	}

	// Verify users were upserted with correct team name
	for _, member := range members {
		user, err := userRepo.GetUser(ctx, member.UserID)
		if err != nil {
			t.Errorf("expected user %s to exist", member.UserID)
		}
		if user.TeamName != teamName {
			t.Errorf("expected user team name %s, got %s", teamName, user.TeamName)
		}
	}
}

func TestTeamService_CreateTeam_AlreadyExists(t *testing.T) {
	teamRepo := newMockTeamRepository()
	userRepo := newMockUserRepository()
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()
	teamName := "backend"
	members := []domain.User{
		{UserID: "user1", Username: "alice", IsActive: true},
	}

	// Create team first time
	_, err := service.CreateTeam(ctx, teamName, members)
	if err != nil {
		t.Fatalf("first create should succeed, got %v", err)
	}

	// Try to create same team again
	_, err = service.CreateTeam(ctx, teamName, members)
	if err == nil {
		t.Fatal("expected error when creating duplicate team")
	}

	if err.Error() != "team already exists" {
		t.Errorf("expected 'team already exists' error, got %v", err)
	}
}

func TestTeamService_CreateTeam_RepositoryError(t *testing.T) {
	teamRepo := newMockTeamRepository()
	teamRepo.createErr = errors.New("database error")
	userRepo := newMockUserRepository()
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()
	teamName := "backend"
	members := []domain.User{
		{UserID: "user1", Username: "alice", IsActive: true},
	}

	_, err := service.CreateTeam(ctx, teamName, members)

	if err == nil {
		t.Fatal("expected error from repository")
	}

	if err.Error() != "database error" {
		t.Errorf("expected 'database error', got %v", err)
	}
}

func TestTeamService_CreateTeam_UpsertUserError(t *testing.T) {
	teamRepo := newMockTeamRepository()
	userRepo := newMockUserRepository()
	userRepo.upsertErr = errors.New("user insert failed")
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()
	teamName := "backend"
	members := []domain.User{
		{UserID: "user1", Username: "alice", IsActive: true},
	}

	_, err := service.CreateTeam(ctx, teamName, members)

	if err == nil {
		t.Fatal("expected error when upserting user")
	}

	if err.Error() != "user insert failed" {
		t.Errorf("expected 'user insert failed', got %v", err)
	}
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	teamRepo := newMockTeamRepository()
	userRepo := newMockUserRepository()
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()
	teamName := "backend"
	members := []domain.User{
		{UserID: "user1", Username: "alice", IsActive: true},
	}

	// Create team first
	createdTeam, _ := service.CreateTeam(ctx, teamName, members)

	// Get team
	team, err := service.GetTeam(ctx, teamName)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if team.TeamName != createdTeam.TeamName {
		t.Errorf("expected team name %s, got %s", createdTeam.TeamName, team.TeamName)
	}
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	teamRepo := newMockTeamRepository()
	userRepo := newMockUserRepository()
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()

	_, err := service.GetTeam(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent team")
	}

	if err.Error() != "team not found" {
		t.Errorf("expected 'team not found', got %v", err)
	}
}

func TestTeamService_GetTeam_RepositoryError(t *testing.T) {
	teamRepo := newMockTeamRepository()
	teamRepo.getErr = errors.New("database error")
	userRepo := newMockUserRepository()
	service := NewTeamService(teamRepo, userRepo)

	ctx := context.Background()

	_, err := service.GetTeam(ctx, "backend")

	if err == nil {
		t.Fatal("expected repository error")
	}

	if err.Error() != "database error" {
		t.Errorf("expected 'database error', got %v", err)
	}
}
