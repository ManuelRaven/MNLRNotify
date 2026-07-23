# Build frontend
FROM node:26-alpine AS builder-frontend
WORKDIR /app

RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile --ignore-scripts

COPY index.html vite.config.ts tsconfig*.json components.d.ts ./
COPY scripts ./scripts
COPY frontend ./frontend
COPY backend ./backend
RUN pnpm run build:client

# Build backend
FROM golang:1.26-alpine AS builder-go
WORKDIR /app

COPY --from=builder-frontend /app/backend .
RUN go mod download
RUN CGO_ENABLED=0 go build -tags production -o mnlrbase

# Deploy binary
FROM alpine:latest AS runner
WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary from GoReleaser
COPY --from=builder-go /app/mnlrbase /mnlrnotify
RUN chmod +x /mnlrnotify

EXPOSE 8090

ENTRYPOINT ["/mnlrnotify", "serve", "--http=0.0.0.0:8090"]
