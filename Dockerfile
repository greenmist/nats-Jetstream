FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o JetStream .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/JetStream .
COPY --from=builder /app/reviews.json .

RUN chmod +x /app/JetStream

CMD ["/app/JetStream"]
