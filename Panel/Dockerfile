FROM registry.docker.ir/golang:1.22.5-bullseye

WORKDIR /app

COPY .env.prod .env

RUN go mod init go-ad-panel

COPY . .

RUN go mod tidy

RUN go build -o panel main.go