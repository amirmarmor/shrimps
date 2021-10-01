FROM node:latest

WORKDIR /usr/src/app
COPY package.json ./
RUN npm update
COPY . ./
RUN npm run build

CMD ["node", "scripts/start.js"]

