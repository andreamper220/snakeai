package handlers

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"snake_ai/internal/server/clients"
	"time"

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

	// TODO add email sending to set is_active=true

	uuid, err := storages.Storage.AddUser(user)
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
	user.Id = uuid

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var userJson shared.UserJson
	if err := json.NewDecoder(r.Body).Decode(&userJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := storages.Storage.GetUserByEmail(userJson.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordMatch, err := user.Password.Check(userJson.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	// TODO move to config session expiration 30sec
	clients.RedisClient.Set(context.Background(), "sessionID_"+session, user.Id.String(), 30*time.Second)
	cookie := http.Cookie{
		Name:     "sessionID",
		Value:    session,
		Path:     "/",
		MaxAge:   30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// TODO move to config secret key 'test'
	// TODO write encrypted cookie to response
}

func UserLogout(w http.ResponseWriter, r *http.Request) {

}
