package handlers_test

import (
	"fmt"
	"net/http"
)

func (s *HandlerTestSuite) TestPlayerMatch() {
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
				http.MethodPost,
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodGet,
			},
			response{
				http.StatusOK,
			},
		},
	}

	s.Register(correctEmail, correctPassword)
	sessionCookie := s.Login(correctEmail, correctPassword)

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /match/", tt.got.method),
			func() {
				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/match", s.Server.URL),
					nil,
				)
				s.Require().NoError(err)
				req.AddCookie(sessionCookie)

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)
				// TODO handle with templates
				s.Equal(tt.want.code, res.StatusCode)
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}
