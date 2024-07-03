package storages

import (
	"github.com/google/uuid"

	matchdata "github.com/andreamper220/snakeai.git/internal/domain/match/data"
	"github.com/andreamper220/snakeai.git/internal/domain/user"
)

var Storage StorageInterface

type StorageInterface interface {
	AddUser(user *user.User) (uuid.UUID, error)
	GetUserByEmail(email string) (*user.User, error)
	IsUserExisted(id uuid.UUID) (bool, error)
	GetPlayerById(id uuid.UUID) (*matchdata.Player, error)
	IncreasePlayerScore(id uuid.UUID) error
}
