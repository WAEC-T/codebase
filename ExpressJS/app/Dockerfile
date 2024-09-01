FROM node:slim

WORKDIR /express-docker

ENV RUN_TYPE=dev

COPY ../package.json ./
COPY ../../.env.local ../

RUN npm install

COPY ../app .

# Exposing server port
EXPOSE 5000

# Running npm run start on port 5000
CMD ["npm", "run", "start"]
