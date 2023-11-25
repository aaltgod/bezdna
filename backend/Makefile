include .env
export

run:
	@sleep 1
	$(info "RUN $(APP_NAME)")
	build
	compose-up
	
	sudo ./bin/$(APP_NAME)

build:
	$(info "BUILD")
	go mod download && go build -o bin/$(APP_NAME) cmd/main.go

clean:
	$(info "CLEAN")
	compose-stop
	rm -rf bin || true

install-goose:
	$(info go install github.com/pressly/goose/v3/cmd/goose@latest)
	go install github.com/pressly/goose/v3/cmd/goose@latest

LOCAL_MIGRATION_DIR=./migrations/postgres
LOCAL_MIGRATION_DSN="host=${POSTGRES_HOST} port=${POSTGRES_PORT} dbname=${POSTGRES_DATABASE} user=${POSTGRES_USERNAME} password=${POSTGRES_PASSWORD} sslmode=disable"

create-new-migration:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres db create ${name} sql

migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

migration-reset:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} reset

compose-up:
	docker compose -p bezdna -f ./docker-compose.yaml up -d
compose-stop:
	docker compose -p bezdna -f ./docker-compose.yaml stop
compose-down:
	docker compose -p bezdna -f ./docker-compose.yaml down
