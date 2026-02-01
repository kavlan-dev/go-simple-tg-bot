FROM golang:latest AS builder
WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/bot

FROM alpine:latest
WORKDIR /app/

COPY --from=builder /app/bot .

CMD ["./bot"]
