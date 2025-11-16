package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jootiee/avito-test-2025/internal/domain"
)

type TeamServiceTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *TeamServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *TeamServiceTestSuite) TestCreate() {
	type args struct {
		teamName string
		members  []domain.User
	}
	type setupResult struct {
		svc      *TeamService
		teamRepo *mockTeamRepository
		userRepo *mockUserRepository
	}

	tests := []struct {
		name       string
		args       args
		setup      func() setupResult
		assertFunc func(t *testing.T, team *domain.Team, err error, sr setupResult)
	}{
		{
			name: "success",
			args: args{
				teamName: "backend",
				members: []domain.User{
					{UserID: "user1", Username: "alice", IsActive: true},
					{UserID: "user2", Username: "bob", IsActive: true},
				},
			},
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				userRepo := newMockUserRepository()
				return setupResult{svc: NewTeamService(teamRepo, userRepo), teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error, sr setupResult) {
				assert.NoError(t, err)
				assert.NotNil(t, team)
				assert.Equal(t, "backend", team.TeamName)
				assert.Len(t, team.Members, 2)
				// Verify users upserted with correct team name
				for _, member := range team.Members {
					user, uerr := sr.userRepo.GetUser(context.Background(), member.UserID)
					assert.NoError(t, uerr)
					assert.Equal(t, "backend", user.TeamName)
				}
			},
		},
		{
			name: "team already exists",
			args: args{teamName: "backend", members: []domain.User{{UserID: "user1", Username: "alice", IsActive: true}}},
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				userRepo := newMockUserRepository()
				svc := NewTeamService(teamRepo, userRepo)
				// Pre-create team
				_, _ = svc.Create(context.Background(), "backend", []domain.User{{UserID: "user1", Username: "alice", IsActive: true}})
				return setupResult{svc: svc, teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error, sr setupResult) {
				assert.Error(t, err)
				assert.Nil(t, team)
				assert.Equal(t, "team already exists", err.Error())
			},
		},
		{
			name: "repository error on create",
			args: args{teamName: "backend", members: []domain.User{{UserID: "user1", Username: "alice", IsActive: true}}},
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				teamRepo.createErr = errors.New("database error")
				userRepo := newMockUserRepository()
				return setupResult{svc: NewTeamService(teamRepo, userRepo), teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error, sr setupResult) {
				assert.Error(t, err)
				assert.Nil(t, team)
				assert.Equal(t, "database error", err.Error())
			},
		},
		{
			name: "upsert user error",
			args: args{teamName: "backend", members: []domain.User{{UserID: "user1", Username: "alice", IsActive: true}}},
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				userRepo := newMockUserRepository()
				userRepo.upsertErr = errors.New("user insert failed")
				return setupResult{svc: NewTeamService(teamRepo, userRepo), teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error, sr setupResult) {
				assert.Error(t, err)
				assert.Nil(t, team)
				assert.Equal(t, "user insert failed", err.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			sr := tt.setup()
			team, err := sr.svc.Create(context.Background(), tt.args.teamName, tt.args.members)
			tt.assertFunc(s.T(), team, err, sr)
		})
	}
}

func (s *TeamServiceTestSuite) TestGet() {
	type setupResult struct {
		svc      *TeamService
		teamRepo *mockTeamRepository
		userRepo *mockUserRepository
	}

	tests := []struct {
		name       string
		teamName   string
		setup      func() setupResult
		assertFunc func(t *testing.T, team *domain.Team, err error)
	}{
		{
			name:     "success",
			teamName: "backend",
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				userRepo := newMockUserRepository()
				svc := NewTeamService(teamRepo, userRepo)
				_, _ = svc.Create(context.Background(), "backend", []domain.User{{UserID: "user1", Username: "alice", IsActive: true}})
				return setupResult{svc: svc, teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, team)
				assert.Equal(t, "backend", team.TeamName)
			},
		},
		{
			name:     "not found",
			teamName: "nonexistent",
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				userRepo := newMockUserRepository()
				return setupResult{svc: NewTeamService(teamRepo, userRepo), teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error) {
				assert.Error(t, err)
				assert.Nil(t, team)
				assert.Equal(t, "team not found", err.Error())
			},
		},
		{
			name:     "repository error",
			teamName: "backend",
			setup: func() setupResult {
				teamRepo := newMockTeamRepository()
				teamRepo.getErr = errors.New("database error")
				userRepo := newMockUserRepository()
				return setupResult{svc: NewTeamService(teamRepo, userRepo), teamRepo: teamRepo, userRepo: userRepo}
			},
			assertFunc: func(t *testing.T, team *domain.Team, err error) {
				assert.Error(t, err)
				assert.Nil(t, team)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			sr := tt.setup()
			team, err := sr.svc.Get(context.Background(), tt.teamName)
			tt.assertFunc(s.T(), team, err)
		})
	}
}

func TestTeamService(t *testing.T) {
	suite.Run(t, new(TeamServiceTestSuite))
}
