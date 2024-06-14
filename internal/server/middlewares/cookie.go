package middlewares

import (
	"encoding/gob"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"snake_ai/internal/server/cookies"
	"snake_ai/internal/server/storages"
)

func WithAuthenticate(h func(w http.ResponseWriter, r *http.Request, userId uuid.UUID), secret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gobEncodedValue, err := cookies.ReadEncrypted(r, "sessionID", secret)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "you are not authorized to access this resource", http.StatusUnauthorized)
			case errors.Is(err, cookies.ErrInvalidValue):
				http.Error(w, "invalid cookie", http.StatusBadRequest)
			default:
				http.Error(w, "something happened getting your cookie data", http.StatusInternalServerError)
			}
			return
		}

		var userID uuid.UUID
		reader := strings.NewReader(gobEncodedValue)
		if err := gob.NewDecoder(reader).Decode(&userID); err != nil {
			http.Error(w, "something happened decoding your cookie data", http.StatusInternalServerError)
			return
		}

		existed, err := storages.Storage.IsUserExisted(userID)
		if err != nil {
			http.Error(w, "something happened getting user from db", http.StatusInternalServerError)
			return
		}
		if !existed {
			http.Error(w, "user not found", http.StatusUnauthorized)
		}

		h(w, r, userID)
	}
}
