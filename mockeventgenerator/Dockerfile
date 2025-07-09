FROM golang:1.22
WORKDIR /
COPY . .

# Remove the replace directive for common and use the local ./common path inside the container
RUN sed -i '/replace common =>/c\replace common => ./common' go.mod
RUN go mod download
RUN go build -o mockeventgenerator

EXPOSE 5672

# Start both RabbitMQ and the mockeventgenerator
CMD ["sh", "-c", "./mockeventgenerator"]
