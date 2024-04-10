# Build stage
FROM golang:1.22-alpine AS builder
LABEL authors="vexsx"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./migration

EXPOSE 80

CMD ["/app/main"]