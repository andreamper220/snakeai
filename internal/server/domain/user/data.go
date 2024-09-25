package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents an object with ID, email and password structure.
type User struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password password  `json:"-"`
}

type password struct {
	Plain *string
	Hash  string
}

func (p *password) Set(plain string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), 14)
	if err != nil {
		return err
	}
	p.Plain = &plain
	p.Hash = string(bytes)
	return nil
}
func (p *password) Check(plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(plain))
	return err == nil, err
}
