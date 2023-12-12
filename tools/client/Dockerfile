FROM alpine

WORKDIR /usr/src/app

COPY . .

RUN apk update && \
    apk add curl

RUN chmod +x client.sh

ENTRYPOINT ["./client.sh"]
