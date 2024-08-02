FROM ubuntu:18.04

RUN apt-get update && apt-get install -y \
    build-essential \
    software-properties-common \
    git \
    curl \
    wget \
    unzip \
    && rm -rf /var/lib/apt/lists/*

RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update && apt-get install -y golang-go

WORKDIR /app

COPY .env.prod .env

RUN go mod init ad-reporter

COPY . .

RUN go mod tidy

RUN go build -o .