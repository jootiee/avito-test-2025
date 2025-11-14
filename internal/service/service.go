package service

// Service aggregates all business logic services
type Service struct {
	Team        *TeamService
	User        *UserService
	PullRequest *PullRequestService
}

// New creates a new Service with all dependencies initialized
func New(
	teamRepo TeamRepository,
	userRepo UserRepository,
	pullRequestRepo PullRequestRepository,
) *Service {
	return &Service{
		Team:        NewTeamService(teamRepo, userRepo),
		User:        NewUserService(userRepo, pullRequestRepo),
		PullRequest: NewPullRequestService(pullRequestRepo, userRepo, teamRepo),
	}
}
