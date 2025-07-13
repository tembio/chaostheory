# Build stage
FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o webserver webserver.go
EXPOSE 8081
CMD ["./webserver"]