build:
	go build -o bin/ ./...

test:
	staticcheck ./...
	go test -v ./...

integrationtest:
	go test ./internal/repository/... -tags integration -v

lint:
	go vet ./...
	golangci-lint run
	sqlc vet
	sqlc verify

run:
	go run ./cmd/server

generate:
	sqlc generate