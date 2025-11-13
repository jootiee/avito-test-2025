package domain

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func New(
	userID, username, teamName string, isActive bool,
) *User {
	return &User{UserID: userID, Username: username, TeamName: teamName, IsActive: isActive}
}

func (u *User) GetID() string {
	return u.UserID
}

func (u *User) GetName() string {
	return u.Username
}

func (u *User) GetIsActive() bool {
	return u.IsActive
}
