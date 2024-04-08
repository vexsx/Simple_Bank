# Build stage
FROM golang:1.22 AS build
LABEL authors="vexsx"

WORKDIR /app

COPY go.* ./

RUN go mod tidy

COPY . .

RUN go build -o main main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration

EXPOSE 80
CMD [ "/app/main" ]
