package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	editorserver "github.com/andreamper220/snakeai/internal/editor/application"
	"github.com/andreamper220/snakeai/internal/server/application"
	"github.com/andreamper220/snakeai/internal/server/domain/user"
)

type HandlerTestSuite struct {
	suite.Suite
	Server       *httptest.Server
	EditorServer *grpc.Server
}

func (s *HandlerTestSuite) SetupTest() {
	editorPort := rand.Intn(50) + 50000
	s.Require().NoError(os.Setenv("ADDRESS", "0.0.0.0:0"))
	s.Require().NoError(os.Setenv("SESSION_SECRET", "1234567887654321"))
	s.Require().NoError(os.Setenv("SESSION_EXPIRATION", "1800"))
	s.Require().NoError(os.Setenv("EDITOR_ADDRESS", "0.0.0.0:"+strconv.Itoa(editorPort)))
	application.ParseFlags()

	s.EditorServer = editorserver.InitGRPCServer(editorPort)
	go func() {
		s.Require().NoError(editorserver.Run(editorPort, true))
	}()
	time.Sleep(500 * time.Millisecond)

	s.Require().NoError(application.Run(true))
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
	s.Require().Equal(http.StatusCreated, res.StatusCode)
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
