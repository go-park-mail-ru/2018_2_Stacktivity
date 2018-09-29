package requests

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Registration struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
	Avatar    string `json:"avatar,omitempty"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}
