# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Run stage
FROM gcr.io/distroless/static

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Run the binary
CMD ["/app/server"]
