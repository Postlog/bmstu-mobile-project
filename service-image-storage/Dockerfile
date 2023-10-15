FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN apk add curl

RUN go mod tidy

RUN go build -v -o ./bin/service ./cmd/service/


ENTRYPOINT ./bin/service
