FROM golang:latest

WORKDIR /application
COPY . .
RUN GOOS=linux GOARCH=amd64 make build
CMD ["./bin/bot"]