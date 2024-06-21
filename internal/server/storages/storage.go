package storages

import (
	"github.com/google/uuid"
	"snake_ai/internal/shared/match/data"

	"snake_ai/internal/shared/user"
)

var Storage StorageInterface

type StorageInterface interface {
	AddUser(user *user.User) (uuid.UUID, error)
	GetUserByEmail(email string) (*user.User, error)
	IsUserExisted(id uuid.UUID) (bool, error)
	GetPlayerById(id uuid.UUID) (*data.Player, error)
	IncreasePlayerScore(id uuid.UUID) error
}
