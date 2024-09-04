.DEFAULT_GOAL := build

check_lint: ## Check if lint tools are installed else install them
	@which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@2023.1.6

check_golang_ci_lint: ## Install 1.55.2 golang-ci-lint version
	@which golangci-lint || echo "Install golangci-lint from https://golangci-lint.run/welcome/install/#local-installation"

.PHONY: check_sqlc
check_sqlc: ## Check if sqlc is installed else install it
	@which sqlc || go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.16.0

.PHONY: check_mockgen
check_mockgen: ## Check if mockgen is installed else install it
	@which mockgen || go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: custom_lint
custom_lint: check_lint check_sqlc ## Check custom lint issues with vet & static check
	# Verify all SQL files
	sqlc compile -f ./internal/db/sqlc.yaml
	# Verify all Go files
	go vet ./...
	staticcheck ./...

.PHONY: lint
lint: check_lint check_sqlc check_golang_ci_lint ## Lint all Go files(using vet and staticcheck)  + SQL files(using sqlc)
	# Verify all SQL files
	sqlc compile -f ./internal/db/sqlc.yaml
	# Verify all Go files
	go vet ./...
	staticcheck ./...
	# Auto fix the linting
	golangci-lint run ./... --fix
	# Verify all Go files with golangci-lint
	golangci-lint run ./...

.PHONY: build
build: ## Build the binary and place it in the `bin` directory. The directory will be created if it does not exist. This is the default target.
# The build tasks delegates to`scripts/build.sh` to build the binary. This is done to avoid the Makefile from becoming too complex.
	@sh -c './scripts/build.sh'

.PHONY: gen_and_build
gen_and_build: gen ## Generate all files and build the binary. Same as `make gen && make build`
	@sh -c './scripts/build.sh'

check_air: ## Check if air is installed else install it. air is used for hot reloading
	which air || go install github.com/cosmtrek/air@latest

.PHONY: local_docker
local_docker: ## Start the required services for the service to run locally (postgres, etc) in docker
	docker compose -f docker-compose-local.yml -p txn-service-local up -d --renew-anon-volumes

# Hot reload for development
.PHONY: dev
dev: check_air check_sqlc local_docker tidy ## Start the service in hot reload mode. This task will start the required docker containers for the service to run locally (postgres, redis, etc) and then start the service in hot reload mode.
	@sh -c './scripts/air.sh' && docker compose -f docker-compose-local.yml -p txn-service-local down

.PHONY: test
test: ## Convenience task for `go test`
	ENVIRONMENT=TEST go test ./...

.PHONY: gen
gen: check_sqlc check_mockgen ## Convenience task for `go generate ./...`
	go generate ./...

.PHONY: race
race: ## Convenience task for `go run -race`
	go run -race ./cmd

.PHONY: check_migrate
check_migrate: ## Check if migrate is installed else install it
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: new_migration
new_migration: check_migrate ## Create a new migration file. Usage: make new_migration name=<migration_name>
	 migrate create -dir=internal/db/migrations/ -seq -ext sql $(name)

.PHONY: local_migrate_one_up
local_migrate_one_up: ## Run a single migration up. Only for local development
	migrate -path='internal/db/migrations' \
			-database='postgres://postgres:password@localhost:5432/pismo?sslmode=disable' \
			up 1

.PHONY: local_migrate_one_down
local_migrate_one_down: ## Run a single migration down. Only for local development
	migrate -path='internal/db/migrations' \
			-database='postgres://postgres:password@localhost:5432/pismo?sslmode=disable' \
			down 1

.PHONY: check_goimports
check_goimports: ## Check if goimports is installed else install it
	which goimports || go install golang.org/x/tools/cmd/goimports@v0.5.0

.PHONY: fmt
fmt: check_goimports ## Format all go files using go fmt + imports
	go fmt ./... && go list ./... | goimports -w .

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: pr
pr: gen fmt tidy lint test ## Run this before raising any Pull Request

.PHONY: help
help: ## Display this help
# This help commands picks up all the comments that start with `##` and prints them in a nice format.
# The comments should be in the following format:<target>:<space><comment>
	@echo "Usage: make <target>"
	@echo "Targets:"
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'
