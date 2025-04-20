# Create Paste Worker Service

This service consumes paste creation messages from RabbitMQ and efficiently inserts them into the database using batch processing.

## Features

- Consumes paste creation messages from RabbitMQ queue
- Uses worker pool pattern with configurable number of goroutines
- Implements batch insertion for improved database performance
- Configurable batch size and flush timeout
- Ensures all pending pastes are inserted on graceful shutdown

## Architecture

```
┌────────────────┐         ┌──────────────┐
│                │         │              │
│   RabbitMQ     │─────────▶ Worker Pool  │
│                │         │              │
└────────────────┘         └──────┬───────┘
                                  │
                                  │ Batch Insert
                                  ▼
                         ┌─────────────────┐
                         │                 │
                         │   Database      │
                         │                 │
                         └─────────────────┘
```

## Running Locally

Run the worker directly:

```
go run main.go
```

## Configuration Options

The worker can be configured via environment variables:

| Environment Variable | Default | Description |
|----------------------|---------|-------------|
| NUM_WORKERS | 5 | Number of worker goroutines |
| BATCH_SIZE | 50 | Number of pastes to collect before performing a batch insert |
| BATCH_TIMEOUT_SECS | 5 | Maximum time to wait before flushing a non-full batch |
| RABBITMQ_URL | amqp://guest:guest@localhost:5672/ | RabbitMQ connection URL |
| DB_HOST | localhost | Database host |
| DB_PORT | 3306 | Database port |
| DB_USER | root | Database username |
| DB_PASSWORD | password | Database password |
| DB_NAME | pastebin | Database name |

## Docker Compose

The worker is designed to run with Docker Compose. The configuration is in the main docker-compose.yml file.

## Batch Processing

The worker implements an intelligent batch processing system to optimize database performance:

1. Pastes are collected in memory until either:
   - The batch size is reached (default: 50 pastes)
   - The batch timeout is reached (default: 5 seconds)

2. When either condition is met, a single SQL statement performs a multi-row insert

3. On service shutdown, any pending pastes in the batch are automatically flushed to the database

This approach significantly reduces database load for high-volume scenarios compared to individual inserts. 