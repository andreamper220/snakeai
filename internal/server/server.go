package server

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"net/http"
	"snake_ai/internal/server/clients"

	"snake_ai/internal/logger"
	"snake_ai/internal/server/handlers"
	"snake_ai/internal/server/storages"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post(`/register`, handlers.UserRegister)
	r.Post(`/login`, handlers.UserLogin)
	r.Post(`/logout`, handlers.UserLogout)

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

func MakeRedis() error {
	opt, err := redis.ParseURL(Config.RedisURL)
	if err != nil {
		return err
	}

	clients.RedisClient = redis.NewClient(opt)
	return nil
}

func Run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	logger.Log.Info("logger established")

	if err := MakeStorage(); err != nil {
		return err
	}
	storage, ok := storages.Storage.(*storages.DBStorage)
	if !ok {
		return errors.New("DB storage is not created")
	}
	defer storage.Connection.Close()
	logger.Log.Info("db connection established")

	if err := MakeRedis(); err != nil {
		return err
	}
	logger.Log.Info("redis connection established")

	logger.Log.Infof("server listening on %s", Config.Address.String())
	return http.ListenAndServe(Config.Address.String(), MakeRouter())
}
