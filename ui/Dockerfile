# syntax=docker/dockerfile:1
FROM node:20-alpine as build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install
COPY . ./
RUN npm run build

FROM nginx:1.24-alpine
COPY --from=build /app/dist /usr/share/nginx/html
