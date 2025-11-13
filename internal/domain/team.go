package domain

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

func NewTeam(name string, members []User) *Team {
	return &Team{TeamName: name, Members: members}
}
