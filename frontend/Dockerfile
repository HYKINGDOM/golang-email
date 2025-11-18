FROM node:18 AS builder
WORKDIR /app
COPY package.json package-lock.json* yarn.lock* pnpm-lock.yaml* ./
RUN npm ci || npm install
COPY . .
RUN npm run build

FROM ubuntu:24.04
WORKDIR /usr/share/nginx/html
RUN apt-get update && apt-get install -y nginx && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/dist /usr/share/nginx/html
COPY public/config.json /usr/share/nginx/html/config.json
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

