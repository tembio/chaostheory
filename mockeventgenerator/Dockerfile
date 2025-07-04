# Build Go app and run RabbitMQ in one container
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN cd mockeventgenerator && go build -o app

# Start both RabbitMQ and the Go app
CMD ["sh", "-c", "./app"]
