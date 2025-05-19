FROM node:22

WORKDIR ./app

COPY ./package*.json .

RUN npm i

COPY ./src .

COPY ./public .

CMD ["npm", "run", "dev"]