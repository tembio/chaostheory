FROM golang:1.22-alpine
ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /
COPY . .

# Remove the replace directive for common and use the local ./common path inside the container
RUN sed -i '/replace common =>/c\replace common => ./common' go.mod
RUN go mod download
RUN go build -o leaderboard .

# Ensure init_db.sh is executable and initialize the database schema
RUN mkdir -p db
RUN chmod +x ./init_db.sh
RUN ./init_db.sh db/leaderboard.db

# Install SQLite runtime
RUN apk add --no-cache sqlite-libs

# Command to run the binary
CMD ["./leaderboard"]
