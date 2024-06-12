package server

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"

	"snake_ai/internal/logger"
)

var Config struct {
	Address     address
	DatabaseDSN string
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

	flag.Var(&addr, "a", "server address host:port")
	flag.StringVar(&Config.DatabaseDSN, "d", "host=snake_db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "database DSN")

	flag.Parse()

	var err error
	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		err = addr.Set(addrEnv)
	}
	if databaseDsnEnv := os.Getenv("DATABASE_DSN"); databaseDsnEnv != "" {
		Config.DatabaseDSN = databaseDsnEnv
	}

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	Config.Address = addr
}
