FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy all files into the container
COPY . .

RUN go build -o JetStream .

FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/JetStream .
COPY --from=builder /app/reviews.json .

RUN chmod +x /app/JetStream

# Run the binary
CMD ["/app/JetStream"]
