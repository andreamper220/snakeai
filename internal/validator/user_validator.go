package validator

import (
	"snake_ai/internal/shared/user"
)

func validateEmail(v *Validator, email string) {
	v.Check(email != "", "email", "email must be provided")
	v.Check(Matches(email, EmailRX), "email", "email is not valid")
}

func validatePlainPassword(v *Validator, password string) {
	v.Check(password != "", "password", "password must be provided")
	v.Check(len(password) >= 8, "password", "password must be more than 8 characters")
	v.Check(len(password) <= 72, "password", "password must be less than 72 characters")
}

func ValidateUser(v *Validator, user *user.User) {
	validateEmail(v, user.Email)
	if user.Password.Plain != nil {
		validatePlainPassword(v, *user.Password.Plain)
	}
}
