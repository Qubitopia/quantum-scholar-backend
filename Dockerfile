FROM golang:latest

WORKDIR /app

# Copy dependency files
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy source code
COPY server/. ./

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

CMD ["./main"]