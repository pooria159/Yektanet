FROM golang:1.22.5-bullseye

WORKDIR /app

COPY .env.prod .env

RUN go mod init ad-reporter

COPY . .

RUN apt-get update && apt-get install -y libc6

RUN go mod tidy

RUN go build -o .