package handlers

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"

	"snake_ai/internal/logger"
	"snake_ai/internal/server/clients"
	"snake_ai/internal/server/cookies"
	"snake_ai/internal/server/storages"
	"snake_ai/internal/shared"
	"snake_ai/internal/validator"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var userJson shared.UserJson
	if err := json.NewDecoder(r.Body).Decode(&userJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &shared.User{
		Email: userJson.Email,
	}
	if err := user.Password.Set(userJson.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := validator.New()
	validator.ValidateUser(v, user)
	if !v.IsValid() {
		http.Error(w, v.String(), http.StatusBadRequest)
		return
	}

	// TODO add email sending OTP to set is_active=true

	userId, err := storages.Storage.AddUser(user)
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
	user.Id = userId

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request, secret []byte, expired time.Duration) {
	var userJson shared.UserJson
	if err := json.NewDecoder(r.Body).Decode(&userJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := storages.Storage.GetUserByEmail(userJson.Email)
	if err != nil {
		switch {
		case errors.Is(err, storages.ErrRecordNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	passwordMatch, err := user.Password.Check(userJson.Password)
	if !passwordMatch {
		http.Error(w, "password does not match", http.StatusUnauthorized)
		return
	}

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&user.Id)
	if err != nil {
		http.Error(w, "something happened encoding your data", http.StatusInternalServerError)
		return
	}
	session := buf.String()

	clients.RedisClient.Set(context.Background(), "sessionID_"+session, user.Id.String(), expired)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("user with ID %s logged in", user.Id)
}

func UserLogout(w http.ResponseWriter, r *http.Request, secret []byte) {
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

	if _, err = clients.RedisClient.Del(context.Background(), "sessionID_"+userID.String()).Result(); err != nil {
		switch {
		case errors.Is(err, redis.Nil):
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

	_, err = w.Write([]byte("You are logged out"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Log.Infof("user with ID %s logged out", userID.String())
}