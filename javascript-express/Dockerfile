FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y curl \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs

WORKDIR /waect-javascript-express

ENV RUN_TYPE=dev

COPY package.json ./

RUN npm install

COPY ./src .

EXPOSE 5000

CMD ["node", "app.js"]