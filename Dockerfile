FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification_service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/notification_service .

ENV POSTGRES_HOST=192.168.1.16
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=admin
ENV POSTGRES_DB=notifications
ENV POSTGRES_SSL_MODE=disable

EXPOSE 8080

CMD ["./notification_service"]
