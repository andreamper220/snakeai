package application

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andreamper220/snakeai/pkg/logger"
)

var Config struct {
	Address             address
	DatabaseDSN         string
	RedisURL            string
	SessionSecret       string
	SessionExpires      time.Duration
	EditorServerAddress address
}

type address struct {
	Host string
	Port int
}

func (a *address) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *address) Set(value string) error {
	var err error
	serverAddress := strings.Split(value, ":")
	if len(serverAddress) != 2 {
		return errors.New("need 2 arguments: host and port")
	}
	a.Host = serverAddress[0]
	a.Port, err = strconv.Atoi(serverAddress[1])

	return err
}

func ParseFlags() {
	addr := address{
		Host: "0.0.0.0",
		Port: 8080,
	}
	editorServerAddr := address{
		Host: "snake-map-editor",
		Port: 50051,
	}
	var sessExpSec int

	if flag.Lookup("a") == nil {
		flag.Var(&addr, "a", "server address host:port")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&Config.DatabaseDSN, "d", "", "database DSN")
	}
	if flag.Lookup("r") == nil {
		flag.StringVar(&Config.RedisURL, "r", "", "redis URL")
	}
	if flag.Lookup("s") == nil {
		flag.StringVar(&Config.SessionSecret, "s", "1234567887654321", "secret to session id encrypt")
	}
	if flag.Lookup("e") == nil {
		flag.IntVar(&sessExpSec, "e", 1800, "session expiration seconds")
	}
	flag.Parse()

	var err error
	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		err = addr.Set(addrEnv)
	}
	if databaseDsnEnv := os.Getenv("DATABASE_DSN"); databaseDsnEnv != "" {
		Config.DatabaseDSN = databaseDsnEnv
	}
	if redisUrlEnv := os.Getenv("REDIS_URL"); redisUrlEnv != "" {
		Config.RedisURL = redisUrlEnv
	}
	if sessionSecretEnv := os.Getenv("SESSION_SECRET"); sessionSecretEnv != "" {
		Config.SessionSecret = sessionSecretEnv
	}
	if sessionExpiresEnv := os.Getenv("SESSION_EXPIRATION"); sessionExpiresEnv != "" {
		sessExpSec, err = strconv.Atoi(sessionExpiresEnv)
	}
	if editorAddrEnv := os.Getenv("EDITOR_ADDRESS"); editorAddrEnv != "" {
		err = editorServerAddr.Set(editorAddrEnv)
	}

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	Config.Address = addr
	Config.SessionExpires = time.Duration(sessExpSec) * time.Second
	Config.EditorServerAddress = editorServerAddr
}
