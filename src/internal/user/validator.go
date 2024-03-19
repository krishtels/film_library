package user

import "film-library/src/internal/tools"

func ValidateCreateUserReuqest(req *CreateUserRequest) *tools.ValidationError {
	ve := &tools.ValidationError{}

	if len(req.Username) == 0 {
		ve.AddViolation("username of length 0")
	}

	if len(req.Password) == 0 {
		ve.AddViolation("password of length 0")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateLoginReuqest(req *LoginRequest) *tools.ValidationError {
	ve := &tools.ValidationError{}

	if len(req.Username) == 0 {
		ve.AddViolation("username of length 0")
	}

	if len(req.Password) == 0 {
		ve.AddViolation("password of length 0")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}
