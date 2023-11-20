FROM golang:1.21.1-alpine AS swaggo
WORKDIR /usr/src/app
RUN go install github.com/swaggo/swag/cmd/swag@latest

FROM golang:1.21.1-alpine AS api
WORKDIR /usr/src/app
COPY . . 

EXPOSE 8080

CMD ["go", "run", "cmd/api/main.go", "-b", "0.0.0.0"]
