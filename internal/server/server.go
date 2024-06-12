package server

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"

	"snake_ai/internal/handlers"
	"snake_ai/internal/logger"
	"snake_ai/internal/server/storages"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post(`/register`, handlers.UserRegister)

	return r
}

func MakeStorage() error {
	if Config.DatabaseDSN == "" {
		return errors.New("database DSN is not set")
	}

	conn, err := sql.Open("pgx", Config.DatabaseDSN)
	if err == nil {
		storages.Storage, err = storages.NewDBStorage(conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func Run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	if err := MakeStorage(); err != nil {
		return err
	}
	storage, ok := storages.Storage.(*storages.DBStorage)
	if !ok {
		return errors.New("DB storage is not created")
	}
	defer storage.Connection.Close()

	logger.Log.Infof("server listening on %s", Config.Address.String())
	return http.ListenAndServe(Config.Address.String(), MakeRouter())
}
