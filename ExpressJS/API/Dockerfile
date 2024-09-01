FROM node:slim

WORKDIR /express-docker

COPY package.json ./
COPY ../../.env.local ../

RUN npm install

COPY ../API .

# Exposing server port
EXPOSE 5001

# Running npm run start on port 5000
CMD ["npm", "run", "start-api"]


