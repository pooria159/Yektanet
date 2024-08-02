FROM golang:1.22.5-bullseye

WORKDIR /app

COPY .env.prod .env

RUN go mod init ad-reporter

COPY . .

RUN go mod tidy

RUN go build -o reporter Reporter.go