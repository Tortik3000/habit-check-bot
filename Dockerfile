FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .

RUN find . -name "*.go" -type f
RUN ls -la cmd/bot/

RUN make build

FROM alpine:latest
#RUN adduser -D -s /bin/sh appuser
WORKDIR /app
COPY --from=builder /app/bin/bot .
#USER appuser
CMD ["./bin/bot"]