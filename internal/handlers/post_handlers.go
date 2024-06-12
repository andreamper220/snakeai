package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
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

	uuid, err := storages.Storage.AddUser(user)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrDuplicateEmail):
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
