# Build stage
FROM golang:1.22 as builder

WORKDIR /app
COPY . .

RUN go build -o main ./cmd

# Final stage
FROM debian:buster-slim

WORKDIR /app

COPY --from=builder /app/main /app/main

CMD ["/app/main"]
