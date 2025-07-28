FROM golang:latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose the port
EXPOSE 8000

# Run the application
CMD ["./main"]
