# Build stage
FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o webserver webserver.go
EXPOSE 8081
EXPOSE 8080
CMD ["./webserver"]