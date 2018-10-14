package responses

import (
	"2018_2_Stacktivity/models"
	"encoding/json"
	"net/http"
	"time"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}

type Validate struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error,omitempty"`
}

type ResponseForm struct {
	ValidateSuccess  bool         `json:"validateSuccess"`
	User             *models.User `json:"user,omitempty"`
	UsernameValidate *Validate    `json:"usernameValidate,omitempty"`
	EmailValidate    *Validate    `json:"emailValidate,omitempty"`
	PasswordValidate *Validate    `json:"passwordValidate,omitempty"`
	Error            *Error       `json:"error,omitempty"`
}

type UserID struct {
	ID int `json:"ID"`
}

func WriteCookie(w http.ResponseWriter, name string, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func Write(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(response)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
