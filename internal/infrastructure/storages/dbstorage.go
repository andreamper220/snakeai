package storages

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"time"

	matchdata "github.com/andreamper220/snakeai/internal/domain/match/data"
	"github.com/andreamper220/snakeai/internal/domain/user"
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
func (dbs *DBStorage) AddUser(user *user.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := dbs.Connection.BeginTx(ctx, nil)
	if err != nil {
		return [16]byte{}, err
	}

	var userID uuid.UUID
	queryUser := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	argsUser := []any{user.Email, user.Password.Hash}
	if err := tx.QueryRowContext(ctx, queryUser, argsUser...).Scan(&userID); err != nil {
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

	queryPlayer := `INSERT INTO players (user_id, name, skill) VALUES ($1, $2, $3)`
	argsPlayer := []any{userID, user.Email, 0}
	if _, err = tx.ExecContext(ctx, queryPlayer, argsPlayer...); err != nil {
		if e := tx.Rollback(); e != nil {
			return [16]byte{}, e
		}

		return [16]byte{}, err
	}

	if err = tx.Commit(); err != nil {
		return [16]byte{}, err
	}

	return userID, nil
}
func (dbs *DBStorage) GetUserByEmail(email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u user.User
	query := `SELECT u.id, u.email, u.password FROM users u WHERE u.email = $1`
	args := []any{email}

	if err := dbs.Connection.QueryRowContext(ctx, query, args...).Scan(
		&u.Id,
		&u.Email,
		&u.Password.Hash,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}
func (dbs *DBStorage) IsUserExisted(id uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	args := []any{id.String()}

	if err := dbs.Connection.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
func (dbs *DBStorage) GetPlayerById(id uuid.UUID) (*matchdata.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	p := matchdata.NewPlayer()
	query := `SELECT p.* FROM players p WHERE p.user_id = $1`
	args := []any{id}

	if err := dbs.Connection.QueryRowContext(ctx, query, args...).Scan(
		&p.Id,
		&p.Name,
		&p.Skill,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &p, nil
}
func (dbs *DBStorage) IncreasePlayerScore(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := dbs.Connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `UPDATE players SET skill = skill + 1 WHERE user_id = $1`
	args := []any{id}
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return err
	}
	return nil
}
