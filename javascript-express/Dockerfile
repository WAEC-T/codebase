FROM node:slim

WORKDIR /express-docker

ENV RUN_TYPE=dev

COPY package.json ./

RUN npm install

COPY ./src .

EXPOSE 5000

CMD ["node", "app.js"]
