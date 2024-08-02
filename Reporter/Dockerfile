FROM golang:1.22.5-bullseye

WORKDIR /app

RUN go mod init reporter

COPY . .

RUN go mod tidy

RUN go build -o reporter Reporter.go