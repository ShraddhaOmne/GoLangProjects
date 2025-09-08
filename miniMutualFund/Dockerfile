# ---------- Build Stage ----------
FROM golang:1.24.5 AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Create app directory and also going into that directory cd /app
WORKDIR /app

# Download dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -ldflags="-s -w" -o app .

# ---------- Final Stage ----------
FROM alpine:latest

# Install SSL certs (needed for HTTPS)
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/app .

EXPOSE 8080

ENTRYPOINT [ "./app" ]
# Run the app
# CMD ["./app"]
