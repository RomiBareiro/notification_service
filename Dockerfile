FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification_service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/notification_service .

EXPOSE 8080

CMD ["./notification_service"]
