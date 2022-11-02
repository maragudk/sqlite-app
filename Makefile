.PHONY: clean
clean:
	rm -f app.db*

.PHONY: cover
cover:
	go tool cover -html=cover.out

.PHONY: lint
lint:
	golangci-lint run

.PHONY: open
open:
	open http://localhost:8080

.PHONY: start
start:
	go run ./cmd/server

.PHONY: test
test:
	go test -coverprofile=cover.out -shuffle on ./...

