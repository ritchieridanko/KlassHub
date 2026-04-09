# ---------- Build Stage ----------
FROM golang:1.25.0-alpine3.22 AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set work directory
WORKDIR /app/services/account

# Copy and download app dependencies
COPY shared ../../shared
COPY services/account/go.mod services/account/go.sum ./
RUN go mod download

# Copy app source
COPY services/account/cmd/app ./cmd/app
COPY services/account/configs ./configs
COPY services/account/internal ./internal

# Build binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.22

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set work directory
WORKDIR /root

# Copy from the Build Stage
COPY --from=builder /app/services/account/bin ./bin
COPY --from=builder /app/services/account/configs ./configs

# Expose port
EXPOSE 50052

# Set entry point
ENTRYPOINT ["./bin/app"]
