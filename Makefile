.PHONY: build
build:
	go build -tags fts5 ./cmd/server

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
	go run -tags fts5 ./cmd/server

.PHONY: test
test:
	go test -tags fts5 -coverprofile=cover.out -shuffle on ./...

