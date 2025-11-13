package entity

type User struct {
	id 		 int
	name 	 string
	isActive bool
}

func New(
	id int,
	name string,
	isActive bool,
) *User {
	return &User{id, name, isActive}
}

func (u *User) GetID() int {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetIsActive() bool {
	return u.isActive
}