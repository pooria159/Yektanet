FROM golang:1.22.5-bullseye

WORKDIR /app

RUN go mod init eventserver

COPY . .

RUN go mod tidy

RUN go build -o eventserver eventserver.go