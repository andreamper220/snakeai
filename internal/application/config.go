package application

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"snake_ai/pkg/logger"
)

var Config struct {
	Address        address
	DatabaseDSN    string
	RedisURL       string
	SessionSecret  string
	SessionExpires time.Duration
}

type address struct {
	host string
	port int
}

func (a *address) String() string {
	return a.host + ":" + strconv.Itoa(a.port)
}

func (a *address) Set(value string) error {
	var err error
	serverAddress := strings.Split(value, ":")
	if len(serverAddress) != 2 {
		return errors.New("need 2 arguments: host and port")
	}
	a.host = serverAddress[0]
	a.port, err = strconv.Atoi(serverAddress[1])

	return err
}

func ParseFlags() {
	addr := address{
		host: "0.0.0.0",
		port: 8080,
	}
	var sessExpSec int

	flag.Var(&addr, "a", "server address host:port")
	flag.StringVar(&Config.DatabaseDSN, "d", "", "database DSN")
	flag.StringVar(&Config.RedisURL, "r", "", "redis URL")
	flag.StringVar(&Config.SessionSecret, "s", "1234567887654321", "secret to session id encrypt")
	flag.IntVar(&sessExpSec, "e", 1800, "session expiration seconds")
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

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	Config.Address = addr
	Config.SessionExpires = time.Duration(sessExpSec) * time.Second
}
