package game_test

import (
	"github.com/andreamper220/snakeai/internal/server/application"
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	"github.com/andreamper220/snakeai/internal/server/domain/user"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/andreamper220/snakeai/pkg/logger"
)

type GameTestSuite struct {
	suite.Suite
	games gamedata.Games
}

func (s *GameTestSuite) SetupTest() {
	if err := logger.Initialize(); err != nil {
		s.Fail(err.Error())
	}
	if err := application.MakeStorage(); err != nil {
		s.Fail(err.Error())
	}

}

func (s *GameTestSuite) AddNewUser(email string) *user.User {
	u := &user.User{
		Email: email,
	}
	if err := u.Password.Set("test_password"); err != nil {
		return nil
	}
	userId, err := storages.Storage.AddUser(u)
	if err != nil {
		return nil
	}
	u.Id = userId

	return u
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
