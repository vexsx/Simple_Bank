# Build stage
FROM golang:1.22.6 AS build
LABEL authors="vexsx"

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/app .

# Final stage
FROM alpine:latest
WORKDIR /app

COPY --from=build /app/app /app/app

EXPOSE 80

CMD ["/app/app"]