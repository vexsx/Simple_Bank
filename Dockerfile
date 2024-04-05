FROM golang:1.22.1
LABEL authors="vexsx"


WORKDIR /app


COPY go.* ./


RUN go mod download


COPY . .


RUN go build -o ./app ./cmd/web/*.go


EXPOSE 80


CMD ["./app"]