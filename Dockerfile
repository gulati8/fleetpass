# Official Golang image (You shouldn't use the `latest` version in production but I'm a bad guy)
FROM golang:latest

# Working directory
WORKDIR /app

# Copy go.mod first (for better caching)
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy everything at /app
COPY . /app

# Build the go app
RUN go build -o main .

# Expose port
EXPOSE 8080

# Define the command to run the app
CMD ["./main"]