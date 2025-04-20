# Get Stats Service

This service provides statistics about pastes in the Pastebin application, with efficient Redis caching.

## Features

- Retrieve monthly statistics for pastes
- Redis caching with different TTLs for current and historical months
- Historical months (not current) are cached longer (7 days)
- Current month data is cached for a shorter period (1 hour)

## Architecture

```
┌────────────────┐         ┌───────────┐
│                │         │           │
│ Get Stats API  │─────────▶   Redis   │ (For non-current months)
│                │         │           │
└────────┬───────┘         └───────────┘
         │
         │ (For current month or cache miss)
         ▼
┌────────────────┐
│                │
│    Database    │
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

### GET /api/stats/monthly?month=YYYY-MM
Retrieves statistics for a specific month.

Example response:
```json
{
  "status": 200,
  "message": "Get monthly stats successfully",
  "data": {
    "totalViews": 1250,
    "avgViewsPerPaste": 8.5,
    "minViews": 0,
    "maxViews": 35
  }
}
```

## Caching Strategy

1. For historical months (not the current month):
   - Check Redis cache first
   - If not found, fetch from database and cache with 7-day TTL
   - Long TTL is appropriate since historical data doesn't change

2. For current month:
   - Always fetch from database to ensure freshness
   - Cache with short 1-hour TTL as data may change frequently

This strategy optimizes database usage while maintaining data accuracy, especially for frequently accessed historical stats. 