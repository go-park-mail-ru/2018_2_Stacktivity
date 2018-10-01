package validate

import (
	"2018_2_Stacktivity/cmd/server-public-api/requests"
	"2018_2_Stacktivity/cmd/server-public-api/responses"
	"regexp"
	"strings"
)

const (
	emailRegexp string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

var (
	rxEmail = regexp.MustCompile(emailRegexp)
)

func LoginValidate(req *requests.Login) (resp *responses.ResponseForm) {
	resp = new(responses.ResponseForm)
	resp.ValidateSuccess = true
	if len(req.Username) == 0 {
		resp.ValidateSuccess = false
		resp.UsernameValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty username"),
		}
	}
	if len(req.Password) == 0 {
		resp.ValidateSuccess = false
		resp.PasswordValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty password"),
		}
	}
	return resp
}

func RegistrationValidate(req *requests.Registration) *responses.ResponseForm {
	resp := new(responses.ResponseForm)
	resp.ValidateSuccess = true
	if len(req.Username) == 0 {
		resp.ValidateSuccess = false
		resp.UsernameValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty username"),
		}
	}
	if len(req.Email) == 0 {
		resp.ValidateSuccess = false
		resp.EmailValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty email"),
		}
	}
	if ok := rxEmail.MatchString(req.Email); !ok {
		resp.ValidateSuccess = false
		resp.PasswordValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Incorrect email"),
		}
	}
	if !strings.EqualFold(req.Password1, req.Password2) {
		resp.ValidateSuccess = false
		resp.PasswordValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Password do not match"),
		}
	}
	if len(req.Password1) == 0 || len(req.Password2) == 0 {
		resp.ValidateSuccess = false
		resp.PasswordValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty password"),
		}
	}
	if !resp.ValidateSuccess {
		resp.Error = responses.NewError("Validate error")
	}
	return resp
}

func UpdateValidate(req *requests.UserUpdate) *responses.ResponseForm {
	resp := new(responses.ResponseForm)
	resp.ValidateSuccess = true
	if len(req.Username) == 0 {
		resp.ValidateSuccess = false
		resp.UsernameValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty username"),
		}
	}
	if len(req.Email) == 0 {
		resp.ValidateSuccess = false
		resp.EmailValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Empty email"),
		}
	}
	if ok := rxEmail.MatchString(req.Email); !ok {
		resp.ValidateSuccess = false
		resp.PasswordValidate = &responses.Validate{
			Success: false,
			Error:   responses.NewError("Incorrect email"),
		}
	}
	if !resp.ValidateSuccess {
		resp.Error = responses.NewError("Validate error")
	}
	return resp
}
