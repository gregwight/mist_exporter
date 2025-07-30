# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mist_exporter ./cmd/main.go

# --- Final Stage ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mist_exporter .
COPY config.yaml.dist ./config.yaml
EXPOSE 9200
CMD ["./mist_exporter", "--config", "config.yaml"]
