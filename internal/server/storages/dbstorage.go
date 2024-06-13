package storages

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"time"

	"snake_ai/internal/shared"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("user not found")
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
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	args := []any{user.Email, user.Password.Hash}

	if err := tx.QueryRowContext(ctx, query, args...).Scan(&userID); err != nil {
		if e := tx.Rollback(); e != nil {
			return [16]byte{}, e
		}

		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return [16]byte{}, ErrDuplicateEmail
		default:
			return [16]byte{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return [16]byte{}, err
	}

	return userID, nil
}

func (dbs *DBStorage) GetUserByEmail(email string) (*shared.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user shared.User
	query := `SELECT u.id, u.email, u.password FROM users u WHERE u.email = $1`
	args := []any{email}

	if err := dbs.Connection.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Email,
		&user.Password.Hash,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
