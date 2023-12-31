COMPOSE_COMMAND = docker compose --env-file configs/.env

.PHONY: all
all: help
help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help\n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Create environment variables
	@chmod +x configs/env.sh && configs/env.sh && mv .env configs/

.PHONY: compose-up
compose-up: ## Run docker compose up for create and start containers
	@${COMPOSE_COMMAND} up -d

.PHONY: compose-build
compose-build: ## Run docker compose up build for create and start containers
	@${COMPOSE_COMMAND} up -d --build

.PHONY: compose-down
compose-down: ## Run docker compose down for stopping and removing containers, networks
	@${COMPOSE_COMMAND} down

.PHONY: compose-remove
compose-remove: ## Run docker compose down for stopping and removing containers, networks, volumes
	@echo -n "All registered data and volumes will be deleted, are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@${COMPOSE_COMMAND} down -v --remove-orphans

.PHONY: compose-exec
compose-exec: ## Run docker compose exec to access bash container
	@${COMPOSE_COMMAND} exec -it backend bash

.PHONY: compose-log
compose-log: ## Run docker compose logs to show logger container
	@${COMPOSE_COMMAND} logs -f backend

.PHONY: compose-top
compose-top: ## Run docker compose top to display the running containers processes
	@${COMPOSE_COMMAND} top

.PHONY: go-fmt
go-fmt: ## Run go fmt
	go fmt ./...

.PHONY: go-vet
go-vet: ## Run go vet
	go vet ./...

.PHONY: go-test
go-test: ## Run go test
	go test ./...

.PHONY: go-test-cover
go-test-cover: ## Run go test with coverage report
	go test -cover ./...

.PHONY: go-test-cover-html
go-test-cover-html: ## Run go test with HTML coverage report
	go test -covermode=count -coverprofile coverage.out -p=1 ./... && \
	go tool cover -html=coverage.out -o coverage.html && \
	xdg-open ./coverage.html
