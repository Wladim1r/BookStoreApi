### Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/apilib ./cmd/main.go

### Final stage
FROM alpine:3.21

WORKDIR /api

COPY --from=builder /app/bin/apilib .
COPY --from=builder /app/internal/database /api/internal/database

EXPOSE 8080

CMD [ "./apilib" ]