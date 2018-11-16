NAME ?= serve
DEV_COMPOSE = dockerfiles/dev/docker-compose.yml
TEST_COMPOSE = dockerfiles/test/docker-compose.yml

.PHONY: help tests db_up db_down tests serve

# Show this help prompt
help:
	@echo '  Usage:'
	@echo ''
	@echo '    make <target>'
	@echo ''
	@echo '  Targets:'
	@echo ''
	@awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':' | grep -v 'IGNORE' | sort | uniq

# Make docker image
image:
	@docker build -t simplinic-task .

# Run golang tests
tests:
	@export CGO_ENABLED=0
	@go test -mod=vendor -v ./...

# Migrate to newest migrations / seeds
db_up:
	@miga --config miga.yml all up

# Roll back all migrations
db_down:
	@miga --config miga.yml migrate down-to 1
	@miga --config miga.yml seed down-to 1

# Run local server
serve:
	@go run -mod=vendor cmd/serve/main.go

.PHONY: ci
# Run tests in docker environment
ci: COMPOSE_FILE=$(TEST_COMPOSE)
ci: NAME=tests
ci: env_run

.PHONY: dev_up dev_down dev_logs dev_deploy dev_restart
# Up dev environment
dev_up: COMPOSE_FILE=$(DEV_COMPOSE)
dev_up: env_up

# Stop and remove dev environment
dev_down: COMPOSE_FILE=$(DEV_COMPOSE)
dev_down: env_down

# Show logs of dev environment
dev_logs: COMPOSE_FILE=$(DEV_COMPOSE)
dev_logs: env_logs

# Deploy $(NAME) service of dev environment
dev_deploy: COMPOSE_FILE=$(DEV_COMPOSE)
dev_deploy: env_deploy

# Restart dev environment
dev_restart: COMPOSE_FILE=$(DEV_COMPOSE)
dev_restart: env_restart

# IGNORE
env_up:
	@time docker-compose -f $(COMPOSE_FILE) up --build -d $(NAME)

# IGNORE
env_run:
	@time docker-compose -f $(COMPOSE_FILE) up --build $(NAME)

# IGNORE
env_restart:
	@time docker-compose -f $(COMPOSE_FILE) restart

# IGNORE
env_deploy:
	@time docker-compose -f $(COMPOSE_FILE) up --build -d $(NAME)

# IGNORE
env_down:
	@time docker-compose -f $(COMPOSE_FILE) down

# IGNORE
env_logs:
	@time docker-compose -f $(COMPOSE_FILE) logs -f --tail 100
