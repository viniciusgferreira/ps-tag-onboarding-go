.PHONY: default run run-with-docs build test docs stop

#Variables
APP_NAME=ps-tag-onboarding-go
COMPOSE_FILE=docker-compose.yml
MAIN_FILE=cmd/ps-tag-onboarding/main.go

# Tasks
default: run-with-docs

run:
	@docker compose -f $(COMPOSE_FILE) up -d
run-with-build:
	@docker compose -f $(COMPOSE_FILE) up -d --build
run-with-docs:
	@swag init -g cmd/ps-tag-onboarding/main.go --output ./docs
	@docker compose -f $(COMPOSE_FILE) up -d
build:
	@go build -o $(APP_NAME) $(MAIN_FILE)
test:
	@go test ./...
docs:
	@swag init -g cmd/ps-tag-onboarding/main.go --output ./docs
stop:
	@docker compose -f $(COMPOSE_FILE) down