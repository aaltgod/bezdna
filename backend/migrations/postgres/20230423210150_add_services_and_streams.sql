-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    IF NOT EXISTS services (
        id BIGSERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        port INT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS streams (
        id BIGSERIAL PRIMARY KEY,
        service_name TEXT NOT NULL,
        service_port INT NOT NULL,
        text TEXT,
        started_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
        ended_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS streams;

DROP TABLE IF EXISTS services;

-- +goose StatementEnd