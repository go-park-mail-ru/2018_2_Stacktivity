package responses

import (
	"2018_2_Stacktivity/cmd/server-public-api/storage"
	"encoding/json"
	"net/http"
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
	ValidateSuccess  bool          `json:"statusSuccess"`
	User             *storage.User `json:"user,omitempty"`
	UsernameValidate *Validate     `json:"usernameValidate,omitempty"`
	EmailValidate    *Validate     `json:"emailValidate,omitempty"`
	PasswordValidate *Validate     `json:"passwordValidate,omitempty"`
	Error            *Error        `json:"error,omitempty"`
}

type Scroreboard struct {
	Board []storage.User `json:"scoreboard"`
}

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(resp)
	if err != nil {
		return err
	}
	return nil
}
