# Multi-stage Dockerfile for Full Stack App

# Stage 1: Build Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./
RUN npm install

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from previous stage
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3: Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=backend-builder /app/main .

# Copy frontend build
COPY --from=backend-builder /app/frontend/build ./frontend/build

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
