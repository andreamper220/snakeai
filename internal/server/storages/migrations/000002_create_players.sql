-- +goose Up
CREATE TABLE IF NOT EXISTS players
(
    user_id UUID NOT NULL PRIMARY KEY REFERENCES users (id),
    name    TEXT NOT NULL UNIQUE,
    skill   INT  NOT NULL
);
CREATE INDEX IF NOT EXISTS users_id_name_indx ON players (user_id, name);

-- +goose Down
DROP TABLE IF EXISTS players;