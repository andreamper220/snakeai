package shared

import (
	"errors"
)

type UserJson struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ErrDuplicateEmail = errors.New("duplicate email")
