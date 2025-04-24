FROM golang:1.22 AS builder

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY backend/ ./

# Install build tools for CGO and SQLite support
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

RUN CGO_ENABLED=1 GOOS=linux go build -o ourchat ./cmd/server


FROM node:23-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build


FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite and SSL
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    sqlite-libs

# Copy the binary from the builder stage
COPY --from=builder /app/backend/ourchat .

# Copy templates and static files from the backend folder
COPY --from=builder /app/backend/templates ./templates
COPY --from=builder /app/backend/static ./static
COPY --from=frontend-builder /app/build ./static-frontend

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./ourchat"]