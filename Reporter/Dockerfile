COPY .env.prod .env

RUN go mod init ad-reporter

COPY . .

RUN go mod tidy

RUN go build -o reporter Reporter.go