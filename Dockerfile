# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# install depedencies
RUN apk add --no-cache curl

# Set working directory in container
WORKDIR /app

# Copy go modules and dependencies into the container at build time
COPY go.mod go.sum ./
RUN go mod download

# Only copy the source code needed for building
COPY . .

# Build the application
# RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service cmd/main/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service cmd/main/main.go

# Stage 2: Create a smaller production image
FROM alpine:latest

# install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata 

# Copy the built binary from builder stage
# COPY --from=builder /app/auth-service . 
COPY --from=builder /app/auth-service /app/auth-service
RUN chmod +x /app/auth-service

# Set working directory 
WORKDIR /app

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/app/auth-service"]
