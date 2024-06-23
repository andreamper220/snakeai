package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	matchdata "snake_ai/internal/domain/match/data"
	matchjson "snake_ai/internal/domain/match/json"
	"strings"
	"time"
)

const (
	correctEmail1    = "test@test.com"
	correctPassword1 = "12345678"
	correctEmail2    = "test@test1.com"
	correctPassword2 = "123456789"
)

func (s *HandlerTestSuite) TestPlayerPartyEnqueue() {
	type request struct {
		method string
		size   int
	}

	type response struct {
		code int
	}

	tests := []struct {
		got  request
		want response
	}{
		{
			request{
				http.MethodGet,
				0,
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodPost,
				1,
			},
			response{
				http.StatusOK,
			},
		},
	}

	s.Register(correctEmail1, correctPassword1)
	sessionCookie1 := s.Login(correctEmail1, correctPassword1)

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /player/party/ %d", tt.got.method, tt.got.size),
			func() {
				pa := matchjson.PartyJson{
					Size:   tt.got.size,
					Width:  20,
					Height: 20,
				}
				body, err := json.Marshal(pa)
				s.Require().NoError(err)
				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/player/party", s.Server.URL),
					bytes.NewBuffer(body),
				)
				s.Require().NoError(err)
				req.AddCookie(sessionCookie1)

				var ws *websocket.Conn
				if tt.want.code == http.StatusOK {
					u := "ws" + strings.TrimPrefix(s.Server.URL, "http") + "/ws"
					header := http.Header{}
					header.Add("Cookie", sessionCookie1.String())
					ws, _, err = websocket.DefaultDialer.Dial(u, header)
					if err != nil {
						s.Fail(err.Error())
					}
				}
				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)

				s.Equal(tt.want.code, res.StatusCode)
				if res.StatusCode == http.StatusOK {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
				out:
					for {
						select {
						case <-ctx.Done():
							if tt.got.size == 1 {
								s.Fail("timeout waiting for websocket message")
							}
							break out
						default:
							_, p, err := ws.ReadMessage()
							s.Require().NoError(err)
							if tt.got.size == 1 {
								var pa matchdata.Party
								s.Require().NoError(json.NewDecoder(bytes.NewBuffer(p)).Decode(&pa))
								break out
							} else {
								s.Fail("unexpected message")
								break out
							}
						}
					}
				}

				if tt.got.size > 0 {
					s.Require().NoError(ws.Close())
				}
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}
