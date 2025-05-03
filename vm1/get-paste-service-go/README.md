# Get Paste Service

This service provides an API to retrieve pastes and asynchronously publishes view updates to RabbitMQ.

## Features

- Retrieve pastes by ID
- Redis caching for improved performance
- Asynchronous view count updates
- RabbitMQ for message publishing
- Direct exchange for message routing

## Architecture

```
┌────────────────┐         ┌───────────┐
│                │         │           │
│ Get Paste API  │─────────▶   Redis   │
│                │         │           │
└────────┬───────┘         └───────────┘
         │
         │                 ┌───────────┐
         │                 │           │
         ├────────────────▶│    DB     │
         │                 │           │
         │                 └───────────┘
         │
         │ Publish
         ▼
┌────────────────┐
│                │
│   RabbitMQ     │────────▶ [get-paste-worker]
│                │
└────────────────┘
```

## Running Locally

Run the service directly:

```
go run main.go
```

## Docker Compose

The service is configured to run with Docker Compose:

```
docker compose up
```

All environment variables are configured in the docker-compose.yml file.

## API Endpoints

### GET /api/paste/:id
Retrieves a paste by its ID.

Example response:
```json
{
  "status": 200,
  "message": "Paste retrieved successfully",
  "data": {
    "id": "abc123",
    "content": "Sample paste content",
    "title": "Sample Paste",
    "language": "plaintext",
    "created_at": "2023-10-15T12:30:45Z",
    "expires_at": "2023-11-15T12:30:45Z",
    "visibility": "public"
  }
}
```

## Caching Strategy

1. On paste retrieval, Redis cache is checked first
2. If not found, data is fetched from the database and cached
3. Cache TTL is set to match paste expiration time or 1 hour by default

## View Update Flow

1. When a paste is successfully retrieved, a view update message is published to RabbitMQ
2. The message is consumed by the separate [get-paste-worker](/get-paste-worker-go) service
3. The worker service handles updating the database with view information 