define find_exe
$(shell for c in $(1); do command -v $$c 2>/dev/null && break; done)
endef

MIGRATE := $(call find_exe, migrate sql-migrate)
check_migrate:
	@if [ -z "$(MIGRATE)"]; then \
		echo "Error: migration tool not found" >&2; exit 1; \
	fi

build:
	go build -o bin/ ./cmd/...

test:
	staticcheck ./...
	go test -v ./...

integrationtest:
	go test ./internal/repository/... -tags integration -v

lint:
	go vet ./...
	golangci-lint run
	sqlc vet

run:
	go run ./cmd/server

generate:
	sqlc generate

migrate-up: check_migrate
	set -a; . ./.env; set +a; $(MIGRATE) -path db/migrations -database "$$DATABASE_URL" up

migrate-down: check_migrate
	set -a; . ./.env; set +a; $(MIGRATE) -path db/migrations -database "$$DATABASE_URL" down $(or $(N),1)