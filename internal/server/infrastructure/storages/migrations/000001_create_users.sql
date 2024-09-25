-- +goose Up
CREATE EXTENSION pgcrypto;
CREATE TABLE IF NOT EXISTS users
(
    id         UUID        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    is_active  BOOLEAN                          DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL             DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS users_id_email_is_active_indx ON users (id, email, is_active);

-- +goose Down
DROP TABLE IF EXISTS users;