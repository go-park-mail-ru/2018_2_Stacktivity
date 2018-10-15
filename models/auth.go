package models

import "gopkg.in/go-playground/validator.v9"

type Login struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

type Registration struct {
	Username string `validate:"required" json:"username"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
	Avatar   string `json:"avatar,omitempty"`
}

func InitValidator() *validator.Validate {
	return validator.New()
}
