FROM golang AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -tags "sqlite_fts5 sqlite_foreign_keys" -ldflags="-s -w" -o /bin/server ./cmd/server

FROM debian:bullseye-slim AS runner
WORKDIR /app

RUN set -x && apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates sqlite3 && \
  rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/server ./

CMD ["./server"]
