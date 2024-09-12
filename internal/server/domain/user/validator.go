package user

import (
	"github.com/andreamper220/snakeai/pkg/validator"
)

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "is not valid")
}

func validatePlainPassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be more than 8 characters")
	v.Check(len(password) <= 72, "password", "must be less than 72 characters")
}

func ValidateUser(v *validator.Validator, user *User) {
	validateEmail(v, user.Email)
	if user.Password.Plain != nil {
		validatePlainPassword(v, *user.Password.Plain)
	}
}
