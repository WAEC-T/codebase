FROM alpine:edge

RUN apk add --no-cache nodejs npm

WORKDIR /waect-javascript-express

ENV RUN_TYPE=dev

COPY package.json ./

RUN npm install

COPY ./src .

EXPOSE 5000

CMD ["node", "app.js"]