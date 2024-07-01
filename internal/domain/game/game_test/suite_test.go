package game_test

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"snake_ai/internal/application"
	gamedata "snake_ai/internal/domain/game/data"
	"snake_ai/internal/domain/user"
	"snake_ai/internal/infrastructure/storages"
	"snake_ai/pkg/logger"
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

func (s *GameTestSuite) AddNewUser() *user.User {
	u := &user.User{
		Email: "test@test.com",
	}
	u.Password.Set("test_password")
	storages.Storage.AddUser(u)

	return u
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
