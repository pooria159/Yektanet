FROM golang:1.22.5-bullseye

WORKDIR /app

RUN go mod init Reporter.go

COPY . .

RUN go mod tidy

RUN go build -o reporter Reporter.go