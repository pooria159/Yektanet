FROM registry.docker.ir/golang:1.22.5-bullseye

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o main main.go