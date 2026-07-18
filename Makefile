.PHONY: bootstrap infra-up infra-down infra-status generate generate-openapi generate-sql db-migrate db-seed db-reset api realtime worker test test-unit test-integration lint verify-backend verify-realtime verify-all web-install web-build web-test android-build android-test

GOOSE_BIN=~/go/bin/goose
SQLC_BIN=~/go/bin/sqlc
OAPI_BIN=~/go/bin/oapi-codegen
DB_URL="postgres://tourtect:change_me_postgres@localhost:5432/tourtect?sslmode=disable"

bootstrap:
	@echo "Installing development tools..."
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Bootstrap completed successfully."

infra-up:
	podman compose up -d postgres redis minio

infra-down:
	podman compose down -v

infra-status:
	podman compose ps

generate-openapi:
	$(OAPI_BIN) -package openapi -generate types,std-http -o backend/generated/openapi/openapi.gen.go backend/api/openapi.yaml

generate-sql:
	cd backend && $(SQLC_BIN) generate

generate: generate-openapi generate-sql

db-migrate:
	$(GOOSE_BIN) -dir backend/db/migrations postgres $(DB_URL) up

db-seed:
	podman exec -i tourtect_postgres_1 psql -U tourtect -d tourtect < backend/db/seed/seed.sql

db-reset:
	$(GOOSE_BIN) -dir backend/db/migrations postgres $(DB_URL) reset

api:
	cd backend && go build -o bin/api cmd/api/main.go

realtime:
	cd backend && go build -o bin/realtime cmd/realtime/main.go

worker:
	cd backend && go build -o bin/worker cmd/worker/main.go

test:
	cd backend && go test ./...

test-unit:
	cd backend && go test -short ./...

test-integration:
	cd backend && go test -run Integration ./...

lint:
	cd backend && go vet ./...

verify-backend:
	@echo "Checking liveness..."
	curl --fail --silent http://localhost:8080/health/live || exit 1
	@echo "Checking readiness..."
	curl --fail --silent http://localhost:8080/health/ready || exit 1
	@echo "Checking places endpoint..."
	curl --fail --silent http://localhost:8080/v1/places || exit 1
	@echo "Checking posts endpoint..."
	curl --fail --silent http://localhost:8080/v1/posts || exit 1
	@echo "Backend verification PASSED"

verify-realtime:
	@echo "Checking realtime health..."
	@# Verify port 8081 is open and accepting requests
	nc -z localhost 8081 || exit 1
	@echo "Realtime verification PASSED"

verify-all: verify-backend verify-realtime

web-install:
	@if [ -d "web" ]; then cd web && npm install; else echo "No web folder found"; fi

web-build:
	@if [ -d "web" ]; then cd web && npm run build; else echo "No web folder found"; fi

web-test:
	@if [ -d "web" ]; then cd web && npm test; else echo "No web folder found"; fi

android-build:
	@if [ -f "android/gradlew" ]; then cd android && ./gradlew assembleDebug; else echo "No Android Gradle project found"; fi

android-test:
	@if [ -f "android/gradlew" ]; then cd android && ./gradlew test; else echo "No Android Gradle project found"; fi
