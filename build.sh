#!/usr/bin/env bash
set -euo pipefail

# ================================
# build-images.sh
# Requirements:
# - Go toolchain installed and on PATH
# - Docker installed and running
# - Run from the Go module root (main package)
# ================================

# Default parameters (can be overridden via env vars or args)
APP_NAME="${APP_NAME:-qs-backend}"
IMAGE_REPO="${IMAGE_REPO:-chetani007/quantum-scholar-backend}"
VERSION_TAG="${VERSION_TAG:-latest}"
CONTEXT_DIR="${CONTEXT_DIR:-.}"
PORT="${PORT:-8000}"
ALPINE_TAG="${ALPINE_TAG:-latest}"

# Output directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUT_DIR="$SCRIPT_DIR/dist"

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

echo "🚀 Building Go binaries for linux/amd64 and linux/arm64..."

# Cross-compile Go binaries
GOOS=linux GOARCH=amd64 go build -o "$OUT_DIR/${APP_NAME}-linux-amd64.out" "$CONTEXT_DIR"
echo "Built $APP_NAME-linux-amd64.out"

GOOS=linux GOARCH=arm64 go build -o "$OUT_DIR/${APP_NAME}-linux-arm64.out" "$CONTEXT_DIR"
echo "Built $APP_NAME-linux-arm64.out"

# Copy mail dir if exists
if [ -d "$SCRIPT_DIR/mail" ]; then
  cp -r "$SCRIPT_DIR/mail" "$OUT_DIR/mail"
fi

# Create Dockerfile directory
DF_DIR="$OUT_DIR/docker"
mkdir -p "$DF_DIR"

# Generate Dockerfiles
cat > "$DF_DIR/Dockerfile.amd64" <<EOF
FROM alpine:${ALPINE_TAG}
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app -u 10000 app
COPY ${APP_NAME}-linux-amd64.out /${APP_NAME}
COPY mail/ /mail/
EXPOSE ${PORT}
USER app
ENTRYPOINT ["/${APP_NAME}"]
EOF

cat > "$DF_DIR/Dockerfile.arm64" <<EOF
FROM alpine:${ALPINE_TAG}
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app -u 10000 app
COPY ${APP_NAME}-linux-arm64.out /${APP_NAME}
COPY mail/ /mail/
EXPOSE ${PORT}
USER app
ENTRYPOINT ["/${APP_NAME}"]
EOF

# Build images
TAG_AMD64="${IMAGE_REPO}:amd64-${VERSION_TAG}"
TAG_ARM64="${IMAGE_REPO}:arm64-${VERSION_TAG}"

echo "🐳 Building Docker image for amd64: $TAG_AMD64"
docker build \
  --file "$DF_DIR/Dockerfile.amd64" \
  --tag "$TAG_AMD64" \
  --platform linux/amd64 \
  "$OUT_DIR"

echo "🐳 Building Docker image for arm64: $TAG_ARM64"
docker build \
  --file "$DF_DIR/Dockerfile.arm64" \
  --tag "$TAG_ARM64" \
  --platform linux/arm64 \
  "$OUT_DIR"

echo "✅ Done. Built images:"
echo " - $TAG_AMD64"
echo " - $TAG_ARM64"
