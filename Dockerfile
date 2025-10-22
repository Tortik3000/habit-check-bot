FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bot ./cmd/bot

FROM alpine:latest
#RUN adduser -D -s /bin/sh appuser
WORKDIR /app
COPY --from=builder /app/bot .
#USER appuser
CMD ["./bot"]