package application

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"net/http"

	"github.com/andreamper220/snakeai/internal/application/handlers/delete_handlers"
	"github.com/andreamper220/snakeai/internal/application/handlers/get_handlers"
	"github.com/andreamper220/snakeai/internal/application/handlers/post_handlers"
	"github.com/andreamper220/snakeai/internal/application/handlers/ws_handlers"
	"github.com/andreamper220/snakeai/internal/application/middlewares"
	gamedata "github.com/andreamper220/snakeai/internal/domain/game/data"
	gameroutines "github.com/andreamper220/snakeai/internal/domain/game/routines"
	matchroutines "github.com/andreamper220/snakeai/internal/domain/match/routines"
	"github.com/andreamper220/snakeai/internal/infrastructure/caches"
	"github.com/andreamper220/snakeai/internal/infrastructure/storages"
	"github.com/andreamper220/snakeai/pkg/logger"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	// w/o auth
	r.Get(`/`, get_handlers.UserAuth)
	r.Post(`/register`, post_handlers.UserRegister)
	r.Post(`/login`, func(w http.ResponseWriter, r *http.Request) {
		post_handlers.UserLogin(w, r, []byte(Config.SessionSecret), Config.SessionExpires)
	})
	// w/ auth
	r.Get(`/match`, middlewares.WithAuthenticate(get_handlers.PlayerMatch, []byte(Config.SessionSecret)))
	r.Post(`/logout`, middlewares.WithAuthenticate(post_handlers.UserLogout, []byte(Config.SessionSecret)))
	r.Post(`/player/party`, middlewares.WithAuthenticate(post_handlers.PlayerPartyEnqueue, []byte(Config.SessionSecret)))
	r.Post(`/player`, middlewares.WithAuthenticate(post_handlers.PlayerEnqueue, []byte(Config.SessionSecret)))
	r.Post(`/player/ai`, middlewares.WithAuthenticate(post_handlers.PlayerRunAi, []byte(Config.SessionSecret)))
	r.Delete(`/player/ai`, middlewares.WithAuthenticate(delete_handlers.PlayerRemoveAi, []byte(Config.SessionSecret)))
	r.Get(`/ws`, middlewares.WithAuthenticate(ws_handlers.PlayerConnection, []byte(Config.SessionSecret)))

	logger.Log.Infof("server listening on %s", Config.Address.String())

	return r
}

func MakeStorage() error {
	if Config.DatabaseDSN != "" {
		conn, err := sql.Open("pgx", Config.DatabaseDSN)
		if err == nil {
			storages.Storage, err = storages.NewDBStorage(conn)
			if err != nil {
				return err
			}
		}
	} else {
		storages.Storage = storages.NewMemStorage()
	}

	return nil
}

func MakeCache() error {
	if Config.RedisURL != "" {
		opt, err := redis.ParseURL(Config.RedisURL)
		if err != nil {
			return err
		}
		caches.Cache = caches.NewRedisCache(opt)
	} else {
		caches.Cache = caches.NewMemCache()
	}

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

	if err := MakeCache(); err != nil {
		return err
	}
	logger.Log.Info("redis connection established")

	numMatchWorkers := 4
	for w := 0; w < numMatchWorkers; w++ {
		go matchroutines.MatchWorker()
	}
	logger.Log.Infof("%d go match workers started", numMatchWorkers)

	numGameWorkers := 8
	gamedata.CurrentGames.Games = make([]*gamedata.Game, 0)
	for w := 0; w < numGameWorkers; w++ {
		go gameroutines.GameWorker()
	}
	logger.Log.Infof("%d go game workers started", numGameWorkers)

	go matchroutines.HandlePartyMessages()
	logger.Log.Info("party messages goroutine started")

	return http.ListenAndServe(Config.Address.String(), MakeRouter())
}
