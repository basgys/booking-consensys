FROM node:16-alpine

WORKDIR /usr/src/web
COPY . .

RUN [ "yarn", "install" ]
CMD [ "yarn", "start" ]