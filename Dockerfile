# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files
COPY --from=builder /app/static ./static

# Copy migrations (for reference, run separately)
COPY --from=builder /app/migrations ./migrations

EXPOSE 3000

CMD ["./server"]
