package models

type User struct {
	ID       int    `db:"uid" json:"id"`
	Username string `db:"username" json:"username" `
	Email    string `db:"email" json:"email"`
	Password string `db:"pass" json:"-"`
	Avatar   string `db:"avatar" json:"avatar,omitempty" `
	Score    int    `db:"score" json:"score"`
}

func NewUser(username string, email string, password string) User {
	return User{
		Username: username,
		Email:    email,
		Password: password,
		Score:    0,
	}
}

type Users []User

func (u Users) Len() int {
	return len(u)
}

func (u Users) Less(i, j int) bool {
	return u[i].Score > u[j].Score
}

func (u Users) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type UserUpdate struct {
	Email    string `json:"email,omitempty" validate:"email"`
	Password string `json:"password,omitempty" valid:"-"`
	Avatar   string `json:"avatar,omitempty" valid:"-"`
}
