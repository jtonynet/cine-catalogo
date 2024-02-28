FROM golang:1.21.1-alpine AS builder
WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY config ./config
COPY internal ./internal

RUN go build ./cmd/api/main.go
RUN chmod a+x main

FROM alpine:latest
WORKDIR /usr/src/app

COPY .env .env
COPY api ./api

COPY --from=builder /usr/src/app/main /usr/src/app/main

CMD ["./main"]