.PHONY: help tests

# Show this help prompt
help:
	@echo '  Usage:'
	@echo ''
	@echo '    make <target>'
	@echo ''
	@echo '  Targets:'
	@echo ''
	@awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':' | grep -v 'IGNORE' | sort | uniq

# Run golang tests
tests:
	@go test -v ./...

# Migrate to newest migrations / seeds
db_up:
	@miga --config miga.yml all up

# Roll back all migrations
db_down:
	@miga --config miga.yml migrate down-to 1
	@miga --config miga.yml seed down-to 1

# Run local server
serve:
	@go run cmd/serve/main.go