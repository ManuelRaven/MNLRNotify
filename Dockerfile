# Build frontend
FROM oven/bun AS builder-bun
WORKDIR /app

COPY . .
RUN bun install --frozen-lockfile
RUN bun run build:client

# Build backend
FROM golang:1.24-alpine AS builder-go
WORKDIR /app

COPY --from=builder-bun /app/backend .
RUN go mod download
RUN CGO_ENABLED=0 go build -tags production -o mnlrbase

# Deploy binary
FROM alpine:latest AS runner

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary from GoReleaser
COPY --from=builder-go /app/mnlrbase /mnlrnotify
RUN chmod +x /mnlrnotify

EXPOSE 8090

ENTRYPOINT ["/mnlrnotify", "serve", "--http=0.0.0.0:8090"]
