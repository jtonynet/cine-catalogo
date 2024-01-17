FROM golang:1.21.1-alpine AS swaggo
WORKDIR /usr/src/app
RUN go install github.com/swaggo/swag/cmd/swag@latest

FROM golang:1.21.1-alpine AS api
WORKDIR /usr/src/app

RUN go build cmd/api/main.go 

COPY ./main main
COPY ./web web
COPY ./api api

EXPOSE 8080

CMD ["./main"]
