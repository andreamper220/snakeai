package application

import (
	"database/sql"
	"github.com/caddyserver/certmagic"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/libdns/cloudflare"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andreamper220/snakeai/internal/server/application/handlers/delete_handlers"
	"github.com/andreamper220/snakeai/internal/server/application/handlers/get_handlers"
	"github.com/andreamper220/snakeai/internal/server/application/handlers/post_handlers"
	"github.com/andreamper220/snakeai/internal/server/application/handlers/ws_handlers"
	"github.com/andreamper220/snakeai/internal/server/application/middlewares"
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	gameroutines "github.com/andreamper220/snakeai/internal/server/domain/game/routines"
	"github.com/andreamper220/snakeai/internal/server/domain/match/routines"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/caches"
	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	"github.com/andreamper220/snakeai/pkg/logger"
	pb "github.com/andreamper220/snakeai/proto"
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
	r.Get(`/editor`, middlewares.WithAuthenticate(get_handlers.PlayerMapEditor, []byte(Config.SessionSecret)))
	r.Post(`/logout`, middlewares.WithAuthenticate(post_handlers.UserLogout, []byte(Config.SessionSecret)))
	r.Post(`/editor/check`, middlewares.WithAuthenticate(post_handlers.PlayerMapCheck, []byte(Config.SessionSecret)))
	r.Post(`/player/party`, middlewares.WithAuthenticate(post_handlers.CreateOrJoinParty, []byte(Config.SessionSecret)))
	r.Post(`/player/party/restore`, middlewares.WithAuthenticate(post_handlers.RestoreParty, []byte(Config.SessionSecret)))
	r.Post(`/player`, middlewares.WithAuthenticate(post_handlers.JoinParty, []byte(Config.SessionSecret)))
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

func ConnectEditor() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(Config.EditorServerAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Run(serverless bool) error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	logger.Log.Info("logger established")

	if err := MakeStorage(); err != nil {
		return err
	}
	storage, ok := storages.Storage.(*storages.DBStorage)
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

	editorConn, err := ConnectEditor()
	if err != nil {
		return err
	}
	grpcclients.EditorClient = pb.NewEditorClient(editorConn)
	logger.Log.Info("grpc editor client connected to server")

	if serverless {
		return nil
	}

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "anrewwolf68@gmail.com"
	certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	dir, _ := filepath.Split(os.Args[0])
	apiTokenFilePath := filepath.Join(dir, "ssl/cloudflare_api_token.txt")
	token, err := os.ReadFile(apiTokenFilePath)
	if err != nil {
		return err
	}
	certmagic.DefaultACME.DNS01Solver = &certmagic.DNS01Solver{
		DNSManager: certmagic.DNSManager{
			DNSProvider: &cloudflare.Provider{
				APIToken: string(token),
			},
		},
	}

	return certmagic.HTTPS([]string{"snakeai.netvolk.online"}, MakeRouter())
}
