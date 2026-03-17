# Pass in the .envrc file, which exports CLMS_DB_DSN
include .envrc

## run: Shortcut to run the API server (same as run/api)
.PHONY: run
run: run/api

## run/api: Run the API server application
.PHONY: run/api
run/api:
	@echo 'Running application…'
	@go run ./cmd/api -port=4000 \
		-env=development \
		-limiter-burst=5 \
		-limiter-rps=2 \
		-limiter-enabled=false \
		-db-dsn=${CLMS_DB_DSN}

## db/psql: Connect to the library database using psql
.PHONY: db/psql
db/psql:
	psql ${CLMS_DB_DSN}

## db/migrations/new name=$1: Create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: Apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${CLMS_DB_DSN} up

## db/migrations/down: Revert all migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Reverting all migrations...'
	migrate -path ./migrations -database ${CLMS_DB_DSN} down

## db/migrations/fix version=$1: Force schema_migrations version
.PHONY: db/migrations/fix
db/migrations/fix:
	@echo 'Forcing schema migrations version to ${version}...'
	migrate -path ./migrations -database ${CLMS_DB_DSN} force ${version}

## db/migrations/init: Create all initial library system migrations
.PHONY: db/migrations/init
db/migrations/init:
	make db/migrations/new name=create_clms_schema