package service

// Service aggregates all business logic services
type Service struct {
	Team *TeamService
	User *UserService
	PR   *PRService
}

// New creates a new Service with all dependencies initialized
func New(teamRepo TeamRepository, userRepo UserRepository, prRepo PRRepository) *Service {
	return &Service{
		Team: NewTeamService(teamRepo, userRepo),
		User: NewUserService(userRepo, prRepo),
		PR:   NewPRService(prRepo, userRepo, teamRepo),
	}
}
