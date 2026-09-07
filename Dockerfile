FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o lottery-app .

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/lottery-app .
COPY --from=builder /app/web ./web
COPY --from=builder /app/bonus ./bonus

EXPOSE 8080

CMD ["./lottery-app"]