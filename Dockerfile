FROM golang:1.22.1
LABEL authors="vexsx"


WORKDIR /app


COPY go.* ./


RUN go mod tidy


COPY . .


RUN go build -o ./app


EXPOSE 80


CMD ["./app"]