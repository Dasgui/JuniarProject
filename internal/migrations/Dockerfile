FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY internal/migrations ./migrations

# goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest