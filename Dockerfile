# Stage 1: Build (Compiler runs on your native architecture)
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Automatically provided by Docker Buildx
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Cross-compile for the target architecture
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o /app/main .

# Stage 2: Final Production Image
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/main /main

# Standard distroless non-root user
USER nonroot:nonroot

EXPOSE 8000
ENTRYPOINT ["/main"]
