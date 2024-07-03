package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andreamper220/snakeai/internal/application"
	gamedata "github.com/andreamper220/snakeai/internal/domain/game/data"
	gameroutines "github.com/andreamper220/snakeai/internal/domain/game/routines"
	matchroutines "github.com/andreamper220/snakeai/internal/domain/match/routines"
	"github.com/andreamper220/snakeai/internal/domain/user"
	"github.com/andreamper220/snakeai/pkg/logger"
)

type HandlerTestSuite struct {
	suite.Suite
	Server *httptest.Server
}

func (s *HandlerTestSuite) SetupTest() {
	application.Config.SessionSecret = "1234567887654321"
	application.Config.SessionExpires = 1800

	if err := logger.Initialize(); err != nil {
		s.Fail(err.Error())
	}
	if err := application.MakeStorage(); err != nil {
		s.Fail(err.Error())
	}
	if err := application.MakeCache(); err != nil {
		s.Fail(err.Error())
	}
	numMatchWorkers := 4
	for w := 0; w < numMatchWorkers; w++ {
		go matchroutines.MatchWorker()
	}
	numGameWorkers := 8
	gamedata.CurrentGames.Games = make([]*gamedata.Game, 0)
	for w := 0; w < numGameWorkers; w++ {
		go gameroutines.GameWorker()
	}
	go matchroutines.HandlePartyMessages()

	s.Server = httptest.NewServer(application.MakeRouter())
}

func (s *HandlerTestSuite) Register(email, password string) uuid.UUID {
	var resU user.User
	u := user.UserJson{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(u)
	s.Require().NoError(err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/register", s.Server.URL),
		bytes.NewBuffer(body),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().NoError(json.NewDecoder(res.Body).Decode(&resU))

	return resU.Id
}

func (s *HandlerTestSuite) Login(email, password string) *http.Cookie {
	u := user.UserJson{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(u)
	s.Require().NoError(err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/login", s.Server.URL),
		bytes.NewBuffer(body),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	s.Require().NoError(err)
	var sessionCookie *http.Cookie = nil
	cookies := res.Cookies()
	for _, c := range cookies {
		if c.Name == "sessionID" {
			sessionCookie = c
		}
	}
	return sessionCookie
}

func (s *HandlerTestSuite) Logout(sessCookie *http.Cookie) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/logout", s.Server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.AddCookie(sessCookie)

	client := &http.Client{}
	_, err = client.Do(req)
	s.Require().NoError(err)
}

func (s *HandlerTestSuite) InitWebSocket(sessCookie *http.Cookie) *websocket.Conn {
	var ws *websocket.Conn
	u := "ws" + strings.TrimPrefix(s.Server.URL, "http") + "/ws"
	header := http.Header{}
	header.Add("Cookie", sessCookie.String())
	ws, _, err := websocket.DefaultDialer.Dial(u, header)
	if err != nil {
		s.Fail(err.Error())
	}

	return ws
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
