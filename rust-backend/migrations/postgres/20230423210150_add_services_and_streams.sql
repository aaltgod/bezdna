-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    IF NOT EXISTS services (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR NOT NULL,
        port INT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS streams (
        id BIGSERIAL,
        service_port INT NOT NULL,
        payload TEXT,
        is_ended BOOLEAN NOT NULL,
        created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS streams;

DROP TABLE IF EXISTS services;

-- +goose StatementEnd