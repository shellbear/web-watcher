FROM golang:1.14 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o web-watcher .

FROM alpine:latest

RUN apk add brotli

WORKDIR /app
COPY --from=builder /app/web-watcher .
ENTRYPOINT ["./web-watcher"]

