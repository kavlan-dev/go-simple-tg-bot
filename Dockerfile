# Сборка
FROM golang AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/bot

# Запуск
FROM alpine
WORKDIR /app/
COPY --from=builder /app/bot .
CMD ["./bot"]
