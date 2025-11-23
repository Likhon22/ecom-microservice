FROM node:18-alpine AS base
WORKDIR /app
COPY package*.json ./


FROM base AS dev
RUN npm install
COPY . .
EXPOSE 5000 5001
CMD [ "npm","run","start:dev" ]

FROM base AS builder
RUN npm install
COPY . .
RUN npm run build


FROM  node:18-alpine AS prod
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY --from=builder /app/dist ./dist
COPY src/proto ./dist/proto
EXPOSE 5000 5001
CMD [ "node","dist/cmd/api/main.js" ]