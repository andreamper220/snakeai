package game_test

import (
	editorserver "github.com/andreamper220/snakeai/internal/editor/application"
	"github.com/andreamper220/snakeai/internal/server/application"
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	"github.com/andreamper220/snakeai/internal/server/domain/user"
	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	pb "github.com/andreamper220/snakeai/proto"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"testing"
	"time"

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

func (s *GameTestSuite) CreateEditorServerWithClient(port int) *grpc.Server {
	editorServer := editorserver.InitGRPCServer()
	go func() {
		s.Require().NoError(editorserver.Run(port, true))
	}()
	time.Sleep(500 * time.Millisecond)

	s.Require().NoError(os.Setenv("EDITOR_ADDRESS", "0.0.0.0:"+strconv.Itoa(port)))
	application.ParseFlags()
	editorConn, err := application.ConnectEditor()
	s.Require().NoError(err)
	grpcclients.EditorClient = pb.NewEditorClient(editorConn)

	return editorServer
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
