# Pass in the .envrc file, which exports CLMS_DB_DSN
-include .envrc

## run: Shortcut to run the API server (same as run/api)
.PHONY: run
run: run/api

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo 'Starting the Community Library API...'
	@ go run ./cmd/api \
		-port=4000 \
		-env=development \
		-limiter-enabled=true \
		-limiter-rps=2 \
		-limiter-burst=5 \
		-cors-trusted-origins="http://localhost:9000 http://localhost:9001" \
		-db-dsn=${CLMS_DB_DSN} \
		-smtp-host=sandbox.smtp.mailtrap.io \
		-smtp-port=2525 \
		-smtp-username=c87b833cae2beb \
		-smtp-password=b4557436bc94d6 \
		-smtp-sender="Community Library <no-reply@bnlsis.org>"

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

## db/seed: reset database and populate with authentic Belizean data
.PHONY: db/seed
db/seed:
	@echo 'Resetting database and injecting seed data...'
	psql "host=localhost dbname=clms user=clms password=clms" -f ./seed.sql

## db/clean: wipe transaction tables for a fresh test run
.PHONY: db/clean
db/clean:
	@echo 'Cleaning transaction tables (Loans/Fines)...'
	psql "host=localhost dbname=clms user=clms password=clms" -c "TRUNCATE Loans, Fine RESTART IDENTITY CASCADE;"

## test/api: test all core logic (including a successful member registration)
.PHONY: test/api
test/api:
	@echo '====================================================='
	@echo '1. TESTING BOOKS: Create a Book'
	@echo '====================================================='
	@curl -s -i -d '{"title": "Automated Test Book", "isbn": "0000000000000", "publisher": "Test Press", "publication_year": 2026, "minimum_age": 5, "description": "Testing the makefile."}' -H "Content-Type: application/json" -X POST http://localhost:4000/v1/books
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '2. TESTING MEMBERS: Registration (Success & Validator)'
	@echo '====================================================='
	@# Test 2a: Successful Registration (to tick the metrics counter)
	@curl -s -i -d '{"first_name": "Tysha", "last_name": "Daniels", "dob": "2006-01-01", "email": "tysha@example.com", "account_status": "Active"}' -H "Content-Type: application/json" -X POST http://localhost:4000/v1/members
	@echo ''
	@echo ''
	@# Test 2b: Validator Rejection
	@curl -s -i -d '{"first_name": "Test", "last_name": "User", "dob": "2000-01-01", "email": "bad-email", "account_status": "Active"}' -H "Content-Type: application/json" -X POST http://localhost:4000/v1/members
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '3. TESTING CIRCULATION: Checkout a Book'
	@echo '====================================================='
	@curl -s -i -d '{"copy_id": 1, "member_id": 1}' -H "Content-Type: application/json" -X POST http://localhost:4000/v1/loans
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '4. TESTING CASHIER: BNLSIS Registration Fee'
	@echo '====================================================='
	@curl -s -i -d '{"member_id": 1, "fine_type": "Local Membership Fee", "amount": 3.00, "paid_status": false}' -H "Content-Type: application/json" -X POST http://localhost:4000/v1/fines
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '--- API AUDIT COMPLETE ---'
	@echo '====================================================='

## test/network: test CORS, metrics, and compression
.PHONY: test/network
test/network:
	@echo '====================================================='
	@echo '1. TESTING CORS (Preflight)'
	@echo '====================================================='
	@curl -s -i -H "Origin: http://localhost:9000" -H "Access-Control-Request-Method: POST" -X OPTIONS http://localhost:4000/v1/books | grep -Ei "Access-Control-Allow-Origin" || true
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '2. TESTING METRICS (Viewing counters)'
	@echo '====================================================='
	@curl -s http://localhost:4000/debug/vars | grep -Ei "total_" || true
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '3. TESTING RATE LIMITER (Bursting)'
	@echo '====================================================='
	@for i in 1 2 3 4 5 6; do \
		curl -s -o /dev/null -w "Request $$i: %{http_code}\n" http://localhost:4000/v1/healthcheck; \
	done
	@echo ''
	@echo ''
	@echo '====================================================='
	@echo '4. TESTING COMPRESSION (GZIP vs Identity)'
	@echo '====================================================='
	@echo 'UNCOMPRESSED (Books):'
	@curl -s -i -H "Accept-Encoding: identity" http://localhost:4000/v1/books | grep -Ei "Content-Length|Content-Encoding|Vary" || true
	@echo ''
	@echo 'COMPRESSED (Metrics - Large File):'
	@# We check for Content-Encoding, Transfer-Encoding, and Vary
	@curl -s -i -H "Accept-Encoding: gzip" http://localhost:4000/debug/vars | grep -aEi "Content-Encoding|Transfer-Encoding|Vary" || true
	@echo ''
	@echo '====================================================='
	@echo '--- INFRASTRUCTURE AUDIT COMPLETE ---'
	@echo '====================================================='