-- +goose Up

-- +goose StatementBegin

CREATE TYPE packet_direction AS ENUM ('IN', 'OUT');

CREATE TABLE
    IF NOT EXISTS services
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    port        INT  NOT NULL UNIQUE,
    flag_regexp TEXT NOT NULL
);

CREATE TABLE
    IF NOT EXISTS packets
(
    id        BIGSERIAL PRIMARY KEY,
    direction packet_direction         NOT NULL,
    payload   TEXT                     NOT NULL,
    stream_id BIGINT                   NOT NULL,
    at        TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX stream_id_idx ON packets (stream_id);

CREATE TABLE
    IF NOT EXISTS streams
(
    id           BIGSERIAL PRIMARY KEY,
    service_port INT NOT NULL
);

CREATE INDEX service_port_idx ON streams (service_port);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS streams;
DROP TABLE IF EXISTS packets;
DROP TABLE IF EXISTS services;
DROP TYPE IF EXISTS packet_direction;

-- +goose StatementEnd