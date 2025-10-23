FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o habit-bot ./cmd/habit-bot

FROM alpine:latest
#RUN adduser -D -s /bin/sh appuser
WORKDIR /app
COPY --from=builder /app/habit-bot .
#USER appuser
CMD ["./bot"]