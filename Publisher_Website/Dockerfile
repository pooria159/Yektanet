FROM registry.docker.ir/golang:1.22.5-bullseye

WORKDIR /app

RUN go mod init go-ad-publisher

COPY . .

RUN go mod tidy

RUN go build -o publisher main.go