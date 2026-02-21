# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /build/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend .
RUN npm run generate

# Stage 2: Build backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /build/backend

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .
RUN CGO_ENABLED=0 go build -o server ./cmd/main.go

# Stage 3: Build whisper
FROM python:3.12-slim AS whisper-builder
RUN pip install --user --no-cache-dir openai-whisper

# Stage 4: Production image
FROM alpine:latest
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    sqlite \
    libc6-compat \
    python3

# Copy whisper from Python builder
COPY --from=whisper-builder /root/.local /root/.local
ENV PATH=/root/.local/bin:$PATH

# Copy built backend
COPY --from=backend-builder /build/backend/server .

# Copy frontend dist to backend public
COPY --from=frontend-builder /build/frontend/.output/public ./public

# Create data directory for database
RUN mkdir -p data tmp

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080 \
    HOST=0.0.0.0 \
    DEVENV=false

# Run server
CMD ["./server"]
