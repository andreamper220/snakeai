package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"snake_ai/internal/domain/user"
)

const (
	correctEmail      = "test@test.com"
	correctPassword   = "12345678"
	incorrectEmail    = "test.com"
	incorrectPassword = "1234"
)

func (s *HandlerTestSuite) TestUserRegister() {
	type request struct {
		method   string
		email    string
		password string
	}

	type response struct {
		email string
		code  int
	}

	tests := []struct {
		got  request
		want response
	}{
		{
			request{
				http.MethodPost,
				incorrectEmail,
				correctPassword,
			},
			response{
				"",
				http.StatusBadRequest,
			},
		},
		{
			request{
				http.MethodPost,
				correctEmail,
				incorrectPassword,
			},
			response{
				"",
				http.StatusBadRequest,
			},
		},
		{
			request{
				http.MethodGet,
				correctEmail,
				correctPassword,
			},
			response{
				"",
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodPost,
				correctEmail,
				correctPassword,
			},
			response{
				correctEmail,
				http.StatusCreated,
			},
		},
	}

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /register/ %s %s", tt.got.method, tt.got.email, tt.got.password),
			func() {
				var resU user.User
				u := user.UserJson{
					Email:    tt.got.email,
					Password: tt.got.password,
				}
				body, err := json.Marshal(u)
				s.Require().NoError(err)

				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/register", s.Server.URL),
					bytes.NewBuffer(body),
				)
				s.Require().NoError(err)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)

				s.Equal(tt.want.code, res.StatusCode)
				if res.StatusCode == http.StatusCreated {
					s.Require().NoError(json.NewDecoder(res.Body).Decode(&resU))
					s.Require().NotEqual(uuid.Nil, resU.Id)
					s.Require().Equal(tt.got.email, resU.Email)
				}

				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestUserLogin() {
	type request struct {
		method   string
		email    string
		password string
	}

	type response struct {
		email string
		code  int
	}

	tests := []struct {
		got  request
		want response
	}{
		{
			request{
				http.MethodPost,
				incorrectEmail,
				correctPassword,
			},
			response{
				"",
				http.StatusNotFound,
			},
		},
		{
			request{
				http.MethodPost,
				correctEmail,
				incorrectPassword,
			},
			response{
				"",
				http.StatusUnauthorized,
			},
		},
		{
			request{
				http.MethodGet,
				correctEmail,
				correctPassword,
			},
			response{
				"",
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodPost,
				correctEmail,
				correctPassword,
			},
			response{
				correctEmail,
				http.StatusOK,
			},
		},
	}

	userId := s.Register(correctEmail, correctPassword)
	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /login/ %s %s", tt.got.method, tt.got.email, tt.got.password),
			func() {
				var resU user.User
				u := user.UserJson{
					Email:    tt.got.email,
					Password: tt.got.password,
				}
				body, err := json.Marshal(u)
				s.Require().NoError(err)

				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/login", s.Server.URL),
					bytes.NewBuffer(body),
				)
				s.Require().NoError(err)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)

				s.Equal(tt.want.code, res.StatusCode)
				if res.StatusCode == http.StatusOK {
					s.Require().NoError(json.NewDecoder(res.Body).Decode(&resU))
					s.Require().Equal(userId, resU.Id)
					s.Require().Equal(tt.got.email, resU.Email)
					var sessionCookie *http.Cookie = nil
					cookies := res.Cookies()
					for _, c := range cookies {
						if c.Name == "sessionID" {
							sessionCookie = c
						}
					}
					s.Require().NotNil(sessionCookie)
				}

				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestUserLogout() {
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

	s.Register(correctEmail, correctPassword)
	sessionCookie := s.Login(correctEmail, correctPassword)

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /logout/", tt.got.method),
			func() {
				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/logout", s.Server.URL),
					nil,
				)
				s.Require().NoError(err)
				req.AddCookie(sessionCookie)

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)

				s.Equal(tt.want.code, res.StatusCode)
				if res.StatusCode == http.StatusOK {
					var sessCookie *http.Cookie = nil
					cookies := res.Cookies()
					for _, c := range cookies {
						if c.Name == "sessionID" {
							sessCookie = c
						}
					}
					s.Equal("", sessCookie.Value)
				}

				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}
