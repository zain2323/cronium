-- +goose Up
ALTER TABLE "users" ADD COLUMN password VARCHAR NOT NULL default '';

-- +goose Down
ALTER TABLE "users" DROP COLUMN password;