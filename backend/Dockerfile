FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o server ./cmd/server

FROM ubuntu:24.04
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/server /app/server
COPY ./config/config.yaml /app/config.yaml
ENV TZ=Asia/Shanghai
ENV EMAIL_DB_HOST=192.168.5.8
ENV EMAIL_DB_PORT=5432
ENV EMAIL_DB_USER=user_S5YD4z
ENV EMAIL_DB_PASSWORD=password_bFkZC5
ENV EMAIL_DB_NAME=email_service
ENV EMAIL_SERVER_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/server"]

