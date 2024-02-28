FROM alpine/git:latest as Git
WORKDIR /hal-explorer
RUN git clone https://github.com/toedter/hal-explorer.git .

FROM node:21.6.2-alpine3.18
WORKDIR /app

COPY --from=Git /hal-explorer/package.json /app
COPY --from=Git /hal-explorer/yarn.lock /app
RUN yarn install --force

COPY --from=Git /hal-explorer /app

EXPOSE 4200

CMD ["yarn", "ng", "serve", "--host=0.0.0.0"]
