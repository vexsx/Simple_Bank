# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY app.env /app/
COPY start.sh /app/
COPY wait-for.sh /app/
COPY db/migration /app/db/migration

# Ensure scripts are executable
RUN chmod +x /app/start.sh /app/wait-for.sh

EXPOSE 8080 9090

# Set entry point and default command
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]