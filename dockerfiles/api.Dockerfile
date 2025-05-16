FROM golang:1.23-alpine AS builder

WORKDIR /app/backend

RUN apk add --no-cache \
    gcc \
    musl-dev \
    pkgconfig \
    sqlite-dev

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=1 GOOS=linux go build -o ourchat ./cmd/server

FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite and SSL
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    sqlite-libs

# Copy the binary from the builder stage
COPY --from=builder /app/backend/ourchat .

# Copy database migrations directory
COPY --from=builder /app/backend/internal/db/migrations ./internal/db/migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./ourchat"]
