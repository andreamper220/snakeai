package storages

import (
	"github.com/google/uuid"
	"snake_ai/internal/shared"
)

var Storage StorageInterface

type StorageInterface interface {
	AddUser(user *shared.User) (uuid.UUID, error)
	GetUserByEmail(email string) (*shared.User, error)
}
