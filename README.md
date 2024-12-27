# MNLRNotify

A focused message routing system with basic Gotify/ntfy.sh compatibility

## About

MNLRNotify implements the core message handling endpoints of Gotify and ntfy.sh, focusing solely on message reception and routing. It provides basic API compatibility for sending messages, but does not support advanced features found in the original services.

## Important Compatibility Notice

This application is **NOT** a full replacement for Gotify or ntfy.sh. It only implements a minimal subset of their APIs:

### Supported Features

- Basic message reception (POST /message endpoint)
- Simple authentication
- Basic priority handling
- Message routing to various services via Shoutrrr

### Unsupported Features

- User management
- Application management
- Websocket connections
- Subscription management
- Client libraries
- Message history
- Web UI for message viewing

If you need these advanced features, please use the original [Gotify](https://gotify.net/) or [ntfy.sh](https://ntfy.sh/) services instead.

## Basic Usage

// ...existing usage documentation...

## Docker Deployment

You can deploy MNLRNotify using the official Docker image from GitHub Container Registry:

```bash
docker run -d \
  --name mnlrnotify \
  -p 8090:8090 \
  -v ./data:/db \
  ghcr.io/manuelraven/mnlrnotify:latest
```

### Docker Compose

For easier deployment, you can use this docker-compose.yml:

```yaml
version: "3"
services:
  mnlrnotify:
    image: ghcr.io/manuelraven/mnlrnotify:latest
    container_name: mnlrnotify
    ports:
      - "8090:8090"
    volumes:
      - ./data:/db
    restart: unless-stopped
```

Start the service with:

```bash
docker-compose up -d
```

### Environment Variables

The Docker container supports the following environment variables:

- `DB_DIR`: Directory for the database files (default: /db)

## API Compatibility

### Gotify-compatible Endpoint

```
POST /message
Content-Type: application/json
X-Gotify-Key: <token>

{
    "message": "Hello World"
}
```
