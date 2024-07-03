package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"time"

	gamejson "snake_ai/internal/domain/game/json"
	matchdata "snake_ai/internal/domain/match/data"
	matchjson "snake_ai/internal/domain/match/json"
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

	userId1 := s.Register(correctEmail1, correctPassword1)
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
								partyUserId := uuid.Nil
								for _, pl := range pa.Players {
									if pl.Id == userId1 {
										partyUserId = pl.Id
										break
									}
								}
								s.Equal(userId1, partyUserId)
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
	s.Logout(sessionCookie1)

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestPlayerEnqueue() {
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
				http.MethodPost,
			},
			response{
				http.StatusOK,
			},
		},
	}

	s.Register(correctEmail1, correctPassword1)
	sessionCookie1 := s.Login(correctEmail1, correctPassword1)
	s.Register(correctEmail2, correctPassword2)
	sessionCookie2 := s.Login(correctEmail2, correctPassword2)

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /player/", tt.got.method),
			func() {
				client := &http.Client{}
				gameIdsChannel := make(chan string, 2)
				if tt.want.code == http.StatusOK {
					// create party
					ws1 := s.InitWebSocket(sessionCookie1)
					go func(ws *websocket.Conn) {
						defer ws.Close()
						for {
							var g1 gamejson.GameJson
							err := ws.ReadJSON(&g1)
							s.Require().NoError(err)
							if g1.Id == "" {
								continue
							}
							gameIdsChannel <- g1.Id
							return
						}
					}(ws1)
					pa := matchjson.PartyJson{
						Size:   2,
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
					// connect to party
					ws2 := s.InitWebSocket(sessionCookie2)
					go func(ws *websocket.Conn) {
						defer ws.Close()
						for {
							var g2 gamejson.GameJson
							err = ws.ReadJSON(&g2)
							s.Require().NoError(err)
							if g2.Id == "" {
								continue
							}
							gameIdsChannel <- g2.Id
							return
						}
					}(ws2)
				}

				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/player", s.Server.URL),
					nil,
				)
				s.Require().NoError(err)
				req.AddCookie(sessionCookie2)
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
							var gameIds [2]string
							n := 0
							for gameId := range gameIdsChannel {
								gameIds[n] = gameId
								n++
								if n == 2 {
									break
								}
							}
							s.Equal(gameIds[0], gameIds[1])
							break out
						}
					}
				}

				s.Require().NoError(res.Body.Close())
			})
	}
	s.Logout(sessionCookie1)
	s.Logout(sessionCookie2)

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestPlayerDelayedEnqueue() {
	s.Register(correctEmail1, correctPassword1)
	sessionCookie1 := s.Login(correctEmail1, correctPassword1)
	s.Register(correctEmail2, correctPassword2)
	sessionCookie2 := s.Login(correctEmail2, correctPassword2)

	s.Run("DELAYED /player/",
		func() {
			client := &http.Client{}
			gameIdsChannel := make(chan string, 2)
			// connect to party
			ws2 := s.InitWebSocket(sessionCookie2)
			go func(ws *websocket.Conn) {
				defer ws.Close()
				for {
					var g2 gamejson.GameJson
					err := ws.ReadJSON(&g2)
					s.Require().NoError(err)
					if g2.Id == "" {
						continue
					}
					gameIdsChannel <- g2.Id
					return
				}
			}(ws2)
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/player", s.Server.URL),
				nil,
			)
			s.Require().NoError(err)
			req.AddCookie(sessionCookie2)
			res, err := client.Do(req)
			s.Require().NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			time.Sleep(3 * time.Second)

			// create party
			ws1 := s.InitWebSocket(sessionCookie1)
			go func(ws *websocket.Conn) {
				defer ws.Close()
				for {
					var g1 gamejson.GameJson
					err := ws.ReadJSON(&g1)
					s.Require().NoError(err)
					if g1.Id == "" {
						continue
					}
					gameIdsChannel <- g1.Id
					return
				}
			}(ws1)
			pa := matchjson.PartyJson{
				Size:   2,
				Width:  20,
				Height: 20,
			}
			body, err := json.Marshal(pa)
			s.Require().NoError(err)
			req, err = http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/player/party", s.Server.URL),
				bytes.NewBuffer(body),
			)
			s.Require().NoError(err)
			req.AddCookie(sessionCookie1)
			_, err = client.Do(req)
			s.Require().NoError(err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
		out:
			for {
				select {
				case <-ctx.Done():
					s.Fail("timeout waiting for websocket message")
					break out
				default:
					var gameIds [2]string
					n := 0
					for gameId := range gameIdsChannel {
						gameIds[n] = gameId
						n++
						if n == 2 {
							break
						}
					}
					s.Equal(gameIds[0], gameIds[1])
					break out
				}
			}

			s.Require().NoError(res.Body.Close())
		})
	s.Logout(sessionCookie1)
	s.Logout(sessionCookie2)

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestPlayerRunAi() {
	type request struct {
		method string
		aiJson gamejson.AiJson
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
				gamejson.AiJson{},
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodPost,
				gamejson.AiJson{
					X:   5,
					Y:   5,
					XTo: 1,
					YTo: 0,
					Ai:  "move,right,left,move,",
				},
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
						for {
							var g1 gamejson.GameJson
							err := ws.ReadJSON(&g1)
							s.Require().NoError(err)
							if len(g1.Snakes.Data) < 1 {
								continue
							}
							for userId := range g1.Snakes.Data {
								userIdsChannel <- userId
								if userId == userId1 {
									return
								}
							}
							s.Fail("no snake for user")
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
				}

				body, err := json.Marshal(tt.got.aiJson)
				s.Require().NoError(err)
				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/player/ai", s.Server.URL),
					bytes.NewBuffer(body),
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
	s.Logout(sessionCookie1)

	defer s.Server.Close()
}
