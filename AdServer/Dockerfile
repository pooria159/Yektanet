FROM registry.docker.ir/golang:1.22.5-bullseye

WORKDIR /app

RUN go mod init go-ad-server

COPY . .

RUN go mod tidy

RUN go build -o adserver main.go