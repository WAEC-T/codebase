FROM alpine:edge

RUN apk add --no-cache curl bash gcc libstdc++ libgcc && curl -fsSL https://bun.sh/install | bash

ENV PATH="/root/.bun/bin:${PATH}"

WORKDIR /waect-javascript-express

ENV RUN_TYPE=dev

COPY package.json ./

RUN bun install

COPY ./src .

EXPOSE 5000

CMD ["bun", "app.js"]
