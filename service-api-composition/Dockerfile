FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN apk add curl

RUN go mod tidy

RUN go build -v -o ./bin/service ./cmd/service/
RUN go build -v -o ./bin/workers/image_scaler ./cmd/workers/image_scaler
