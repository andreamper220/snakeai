package application

import (
	"database/sql"
	"github.com/andreamper220/snakeai/internal/server/application/handlers/delete_handlers"
	get_handlers2 "github.com/andreamper220/snakeai/internal/server/application/handlers/get_handlers"
	post_handlers2 "github.com/andreamper220/snakeai/internal/server/application/handlers/post_handlers"
	"github.com/andreamper220/snakeai/internal/server/application/handlers/ws_handlers"
	"github.com/andreamper220/snakeai/internal/server/application/middlewares"
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	gameroutines "github.com/andreamper220/snakeai/internal/server/domain/game/routines"
	"github.com/andreamper220/snakeai/internal/server/domain/match/routines"
	caches2 "github.com/andreamper220/snakeai/internal/server/infrastructure/caches"
	storages2 "github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"net/http"

	"github.com/andreamper220/snakeai/pkg/logger"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	// w/o auth
	r.Get(`/`, get_handlers2.UserAuth)
	r.Post(`/register`, post_handlers2.UserRegister)
	r.Post(`/login`, func(w http.ResponseWriter, r *http.Request) {
		post_handlers2.UserLogin(w, r, []byte(Config.SessionSecret), Config.SessionExpires)
	})
	// w/ auth
	r.Get(`/match`, middlewares.WithAuthenticate(get_handlers2.PlayerMatch, []byte(Config.SessionSecret)))
	r.Post(`/logout`, middlewares.WithAuthenticate(post_handlers2.UserLogout, []byte(Config.SessionSecret)))
	r.Post(`/player/party`, middlewares.WithAuthenticate(post_handlers2.PlayerPartyEnqueue, []byte(Config.SessionSecret)))
	r.Post(`/player`, middlewares.WithAuthenticate(post_handlers2.PlayerEnqueue, []byte(Config.SessionSecret)))
	r.Post(`/player/ai`, middlewares.WithAuthenticate(post_handlers2.PlayerRunAi, []byte(Config.SessionSecret)))
	r.Delete(`/player/ai`, middlewares.WithAuthenticate(delete_handlers.PlayerRemoveAi, []byte(Config.SessionSecret)))
	r.Get(`/ws`, middlewares.WithAuthenticate(ws_handlers.PlayerConnection, []byte(Config.SessionSecret)))

	logger.Log.Infof("server listening on %s", Config.Address.String())

	return r
}

func MakeStorage() error {
	if Config.DatabaseDSN != "" {
		conn, err := sql.Open("pgx", Config.DatabaseDSN)
		if err == nil {
			storages2.Storage, err = storages2.NewDBStorage(conn)
			if err != nil {
				return err
			}
		}
	} else {
		storages2.Storage = storages2.NewMemStorage()
	}

	return nil
}

func MakeCache() error {
	if Config.RedisURL != "" {
		opt, err := redis.ParseURL(Config.RedisURL)
		if err != nil {
			return err
		}
		caches2.Cache = caches2.NewRedisCache(opt)
	} else {
		caches2.Cache = caches2.NewMemCache()
	}

	return nil
}

func Run(serverless bool) error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	logger.Log.Info("logger established")

	if err := MakeStorage(); err != nil {
		return err
	}
	storage, ok := storages2.Storage.(*storages2.DBStorage)
	if ok {
		defer storage.Connection.Close()
		logger.Log.Info("db connection established")
	}

	if err := MakeCache(); err != nil {
		return err
	}
	logger.Log.Info("cache established")

	numMatchWorkers := 4
	for w := 0; w < numMatchWorkers; w++ {
		go match_routines.MatchWorker()
	}
	logger.Log.Infof("%d go match workers started", numMatchWorkers)

	numGameWorkers := 8
	gamedata.CurrentGames.Games = make([]*gamedata.Game, 0)
	for w := 0; w < numGameWorkers; w++ {
		go gameroutines.GameWorker()
	}
	logger.Log.Infof("%d go game workers started", numGameWorkers)

	go match_routines.HandlePartyMessages()
	logger.Log.Info("party messages goroutine started")

	if serverless {
		return nil
	}
	return http.ListenAndServe(Config.Address.String(), MakeRouter())
}
