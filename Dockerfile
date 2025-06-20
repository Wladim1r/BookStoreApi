### Build stage
FROM golang:1.24.2 AS builder

# Устанавливаем зависимости для librdkafka
RUN apt-get update && \
    apt-get install -y gcc librdkafka-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/bin/apilib ./cmd/main.go

### Final stage
FROM debian:bookworm-slim

WORKDIR /api

COPY --from=builder /app/bin/apilib .
COPY --from=builder /app/internal/database /api/internal/database

EXPOSE 8080

CMD ["./apilib"]
