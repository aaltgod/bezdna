-- +goose Up

-- +goose StatementBegin

CREATE TYPE flag_direction AS ENUM ('IN', 'OUT');

CREATE TABLE
    IF NOT EXISTS services
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    port        INT  NOT NULL UNIQUE,
    flag_regexp TEXT NOT NULL
);

CREATE TABLE
    IF NOT EXISTS streams
(
    id           BIGSERIAL PRIMARY KEY,
    service_port INT                      NOT NULL,
    payload      TEXT                     NOT NULL,
    started_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at     TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX service_port_idx ON streams (service_port);

CREATE TABLE
    IF NOT EXISTS flags
(
    id        BIGSERIAL PRIMARY KEY,
    text      TEXT           NOT NULL,
    regexp    TEXT           NOT NULL,
    stream_id BIGINT         NOT NULL,
    direction flag_direction NOT NULL
);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS flags;
DROP TABLE IF EXISTS streams;
DROP TABLE IF EXISTS services;
DROP TYPE IF EXISTS flag_direction;

-- +goose StatementEnd