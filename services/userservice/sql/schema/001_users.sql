-- +goose Up
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    first_name VARCHAR   NOT NULL,
    last_name  VARCHAR   NOT NULL,
    email      VARCHAR NOT NULL UNIQUE ,
    phone      VARCHAR NOT NULL UNIQUE ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;