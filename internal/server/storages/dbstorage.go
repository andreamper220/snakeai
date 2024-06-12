package storages

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DBStorage struct {
	Connection *sql.DB
}

func NewDBStorage(conn *sql.DB) (*DBStorage, error) {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return nil, err
	}
	if err := goose.Up(conn, "migrations"); err != nil {
		return nil, err
	}

	return &DBStorage{
		Connection: conn,
	}, nil
}
func (dbs *DBStorage) AddUser(user User) error {
	return nil
}
