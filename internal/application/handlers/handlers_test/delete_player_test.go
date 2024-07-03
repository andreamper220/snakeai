package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	gamejson "snakeai/internal/domain/game/json"
	matchjson "snakeai/internal/domain/match/json"
)

func (s *HandlerTestSuite) TestPlayerRemoveAi() {
	type request struct {
		method string
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
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodDelete,
			},
			response{
				http.StatusOK,
			},
		},
	}

	userId1 := s.Register(correctEmail1, correctPassword1)
	sessionCookie1 := s.Login(correctEmail1, correctPassword1)

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /player/ai/", tt.got.method),
			func() {
				userIdsChannel := make(chan uuid.UUID, 1)
				client := &http.Client{}
				if tt.want.code == http.StatusOK {
					// create party
					ws1 := s.InitWebSocket(sessionCookie1)
					go func(ws *websocket.Conn) {
						defer ws.Close()
						userId := uuid.Nil
						for {
							var g1 gamejson.GameJson
							err := ws.ReadJSON(&g1)
							s.Require().NoError(err)
							if g1.Id != "" {
								if len(g1.Snakes.Data) == 0 && userId == uuid.Nil {
									continue
								}
								if len(g1.Snakes.Data) > 0 {
									for snakeUserId := range g1.Snakes.Data {
										userId = snakeUserId
									}
									continue
								}
								userIdsChannel <- userId
								return
							} else {
								continue
							}
						}
					}(ws1)
					pa := matchjson.PartyJson{
						Size:   1,
						Width:  20,
						Height: 20,
					}
					body, err := json.Marshal(pa)
					s.Require().NoError(err)
					req, err := http.NewRequest(
						http.MethodPost,
						fmt.Sprintf("%s/player/party", s.Server.URL),
						bytes.NewBuffer(body),
					)
					s.Require().NoError(err)
					req.AddCookie(sessionCookie1)
					_, err = client.Do(req)
					s.Require().NoError(err)

					// add snake
					body, err = json.Marshal(gamejson.AiJson{
						X:   5,
						Y:   5,
						XTo: 1,
						YTo: 0,
						Ai:  "move,move,move,",
					})
					s.Require().NoError(err)
					req, err = http.NewRequest(
						http.MethodPost,
						fmt.Sprintf("%s/player/ai", s.Server.URL),
						bytes.NewBuffer(body),
					)
					s.Require().NoError(err)
					req.AddCookie(sessionCookie1)
					_, err = client.Do(req)
					s.Require().NoError(err)
				}

				time.Sleep(1 * time.Second)
				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/player/ai", s.Server.URL),
					nil,
				)
				s.Require().NoError(err)
				req.AddCookie(sessionCookie1)
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
							s.Fail("timeout waiting for websocket message")
							break out
						default:
							userId := <-userIdsChannel
							s.Equal(userId1, userId)
							break out
						}
					}
				}

				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}
