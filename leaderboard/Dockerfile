# Start from the official Golang image
FROM golang:1.22-alpine as builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Remove the replace directive for common and use the local ./common path inside the container
RUN sed -i '/replace common =>/c\replace common => ./common' go.mod

# Download dependencies (including SQLite driver)
RUN go mod download

# Build the Go app (static binary)
RUN go build -o leaderboard .

# Final image
FROM alpine:latest

# Install SQLite runtime
RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/leaderboard .

# Expose a port if your app serves HTTP (optional)
# EXPOSE 8080

# Command to run the binary
CMD ["./leaderboard"]
