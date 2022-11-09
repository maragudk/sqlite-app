.PHONY: build
build:
	go build -tags "sqlite_fts5 sqlite_foreign_keys" ./cmd/server

.PHONY: clean
clean:
	rm -f app.db*

.PHONY: cover
cover:
	go tool cover -html=cover.out

.PHONY: lint
lint:
	golangci-lint run

.PHONY: migrate-down
migrate-down:
	go run ./cmd/migrate down

.PHONY: migrate-up
migrate-up:
	go run ./cmd/migrate up

.PHONY: open
open:
	open http://localhost:8080

.PHONY: start
start:
	go run -tags "sqlite_fts5 sqlite_foreign_keys" ./cmd/server

.PHONY: test
test:
	go test -tags "sqlite_fts5 sqlite_foreign_keys" -coverprofile=cover.out -shuffle on ./...

