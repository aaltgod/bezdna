APP_NAME=backend

run: build
	$(info "RUN docker containers")
	docker compose up -d
	@sleep 1
	$(info "RUN $(APP_NAME)")
	sudo ./bin/$(APP_NAME)

build: clean
	$(info "BUILD")
	go mod download && go build -o bin/$(APP_NAME) cmd/main.go

clean:
	$(info "STOP docker containers")
	docker compose stop
	$(info "CLEAN")	
	rm -rf bin || true

install-goose:
	$(info go install github.com/pressly/goose/v3/cmd/goose@latest)
	go install github.com/pressly/goose/v3/cmd/goose@latest

LOCAL_MIGRATION_DIR=./migrations/postgres
LOCAL_MIGRATION_DSN="host=localhost port=5432 dbname=bezdna user=user password=1234 sslmode=disable"

create-new-migration:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres db create ${name} sql

migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

migration-reset:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} reset