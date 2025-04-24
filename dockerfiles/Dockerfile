FROM golang:1.23-alpine AS builder

WORKDIR /app/backend

RUN apk add --no-cache \
    gcc \
    musl-dev \
    pkgconfig \
    sqlite-dev \

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

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