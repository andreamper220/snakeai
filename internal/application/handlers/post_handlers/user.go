package post_handlers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/andreamper220/snakeai/internal/infrastructure/caches"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/andreamper220/snakeai/internal/application/cookies"
	gamedata "github.com/andreamper220/snakeai/internal/domain/game/data"
	"github.com/andreamper220/snakeai/internal/domain/user"
	"github.com/andreamper220/snakeai/internal/domain/ws"
	"github.com/andreamper220/snakeai/internal/infrastructure/storages"
	"github.com/andreamper220/snakeai/pkg/logger"
	"github.com/andreamper220/snakeai/pkg/validator"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var userJson user.UserJson
	if err := json.NewDecoder(r.Body).Decode(&userJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := &user.User{
		Email: userJson.Email,
	}
	if err := u.Password.Set(userJson.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := validator.New()
	user.ValidateUser(v, u)
	if !v.IsValid() {
		http.Error(w, v.String(), http.StatusBadRequest)
		return
	}

	// TODO add email sending OTP to set is_active=true

	userId, err := storages.Storage.AddUser(u)
	if err != nil {
		switch {
		case errors.Is(err, storages.ErrDuplicateEmail):
			v.AddError("email", "user with this email already exists")
			http.Error(w, v.String(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	u.Id = userId

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request, secret []byte, expired time.Duration) {
	var userJson user.UserJson
	if err := json.NewDecoder(r.Body).Decode(&userJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := storages.Storage.GetUserByEmail(userJson.Email)
	if err != nil {
		switch {
		case errors.Is(err, storages.ErrRecordNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	passwordMatch, err := u.Password.Check(userJson.Password)
	if !passwordMatch {
		http.Error(w, "password does not match", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "something happened decoding your password", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&u.Id)
	if err != nil {
		http.Error(w, "something happened encoding your data", http.StatusInternalServerError)
		return
	}
	session := buf.String()

	if err := caches.Cache.AddSession(session, u.Id.String(), expired); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, "something happened setting your cache data", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "sessionID",
		Value:    session,
		Path:     "/",
		MaxAge:   int(expired.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	if err := cookies.WriteEncrypted(w, cookie, secret); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, "something happened setting your cookie data", http.StatusInternalServerError)
		return
	}

	ws.Connections.Remove(u.Id)
	gamedata.RemovePlayer(u.Id)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("user with ID %s logged in", u.Id)
}

func UserLogout(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(&userId)
	if err != nil {
		http.Error(w, "something happened checking your user id", http.StatusInternalServerError)
		return
	}

	if err = caches.Cache.DelSession(buf.String()); err != nil {
		switch {
		case errors.Is(err, caches.ErrNoSession):
			http.Error(w, "you are not authorized to access this resource", http.StatusUnauthorized)
		default:
			http.Error(w, "something happened deleting your cache data", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Now(),
	})

	ws.Connections.Remove(userId)
	gamedata.RemovePlayer(userId)

	_, err = w.Write([]byte("You are logged out"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Log.Infof("user with ID %s logged out", userId.String())
}
