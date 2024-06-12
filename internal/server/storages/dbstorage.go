package storages

import (
	"context"
	"database/sql"
	"embed"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"time"

	"snake_ai/internal/shared"
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
func (dbs *DBStorage) AddUser(user *shared.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := dbs.Connection.BeginTx(ctx, nil)
	if err != nil {
		return [16]byte{}, err
	}

	var userID uuid.UUID
	queryUser := `
    INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	argsUser := []any{user.Email, user.Password.Hash}

	if err := tx.QueryRowContext(ctx, queryUser, argsUser...).Scan(&userID); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return [16]byte{}, shared.ErrDuplicateEmail
		default:
			return [16]byte{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return [16]byte{}, err
	}

	return userID, nil
}
