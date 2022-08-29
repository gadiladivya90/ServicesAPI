# make: app info
APP_NAME     := konnectServices
APP_WORKDIR  := $(shell pwd)
BUILD_DIR := $(APP_WORKDIR)/build
APP_PACKAGES := $(shell go list -f '{{.Dir}}' ./...)
APP_LOG_FMT  := `/bin/date "+%Y-%m-%d %H:%M:%S %z [$(APP_NAME)]"`
DATABASE_DOCKER_COMPOSE := $(APP_WORKDIR)/storage/docker

.PHONY: up
up: build-binary app db home run

.PHONY: down
down: app-down db-down home

.PHONY: restart
restart: down up 

.PHONY: run
run: 
	@SERVER_ADDRESS=localhost SERVER_PORT=9090 \
	DB_HOST=localhost DB_PORT=5432 DB_USERNAME=postgres DB_PASSWORD=postgres DB_NAME=services \
		nohup go run main.go &  

.PHONY: status
status: 
	@docker compose ps \
		&& docker compose logs api

.PHONY: app
app: build-binary
	@docker compose up --build --remove-orphans --detach 

.PHONY: db
db:
	@cd ${DATABASE_DOCKER_COMPOSE} \
		&& docker compose up --build --remove-orphans --detach

.PHONY: app-down
app-down:
	@docker compose down -v --rmi local --remove-orphans

.PHONY: db-down
db-down:
	@cd ${DATABASE_DOCKER_COMPOSE} \
		&& docker compose down -v --rmi local --remove-orphans


.PHONY: home
home:
	@cd ${APP_WORKDIR}



.PHONY: build-clean
build-clean:
	@rm -rf $(BUILD_DIR)

.PHONY: build-binary
build-binary: build-clean 
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build \
		-o $(BUILD_DIR) -ldflags '-extldflags "-static"' \
		main.go

.PHONY: build-generate
build-generate: 
	@SWAGGER_GENERATE_EXTENSION=false go generate ./...

.PHONY:tests
tests: generate-mocks run-tests
	
.PHONY:run-tests
run-tests:
	@./run_tests.sh

.PHONY:genrate-mocks
run-tests:
	@./generate_mocks.sh	