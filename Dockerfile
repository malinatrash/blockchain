FROM golang:1.22 as builder

WORKDIR /app
COPY . .

RUN go build -o main .

FROM debian:buster-slim

COPY --from=builder /app/main /app/main

CMD ["/app/main"]
