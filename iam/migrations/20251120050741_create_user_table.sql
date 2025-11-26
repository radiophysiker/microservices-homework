-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    uuid                UUID PRIMARY KEY,
    login               TEXT NOT NULL,
    email               TEXT NOT NULL,
    password_hash       TEXT NOT NULL,
    notification_methods JSONB NOT NULL DEFAULT '[]'::jsonb,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex
    ON users (login);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex
    ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
