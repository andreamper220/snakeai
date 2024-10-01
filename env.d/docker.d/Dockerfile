FROM golang:1.22.1-bookworm as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o server ./cmd/server
RUN go build -v -o editor ./cmd/editor

FROM debian:bookworm-slim as server
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/ssl/fullchain.pem /app/internal/server/application/ssl/fullchain.pem
COPY --from=builder /app/ssl/privkey.pem /app/internal/server/application/ssl/privkey.pem
COPY --from=builder /app/internal/server/templates /app/internal/server/templates
COPY --from=builder /app/server /app/server

CMD ["/app/server", \
     "-d", "host=snake_db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", \
     "-r", "redis://snake_redis:6379" \
     ]

FROM debian:bookworm-slim as editor
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/editor /app/editor

CMD ["/app/editor"]