-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS services (
    name VARCHAR (30) NOT NULL,
    port INT NOT NULL,
    PRIMARY KEY(name, port)  
);

CREATE TABLE IF NOT EXISTS streams (
    service_name VARCHAR (30) NOT NULL,
    service_port INT NOT NULL,
    ack BIGSERIAL NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    payload TEXT,
    FOREIGN KEY(service_name, service_port)
        REFERENCES services(name, port)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS streams;
DROP TABLE IF EXISTS services;
-- +goose StatementEnd
