# Get Paste Worker

This is a standalone worker service that processes view updates for the pastebin application. It consumes messages from RabbitMQ and updates the database with view counts.

## Features

- Consumes view update messages from RabbitMQ
- Uses goroutine pool for efficient processing
- Updates paste view records in the database
- Graceful shutdown handling
- Configurable number of worker goroutines via environment variable

## Architecture

```
┌────────────────┐
│                │
│   RabbitMQ     │
│                │
└────────┬───────┘
         │
         │ Consume 
         ▼                
┌────────────────┐        ┌───────────┐
│                │        │           │
│ Get Paste      │───────▶│    DB     │
│ Worker Pool    │        │           │
└────────────────┘        └───────────┘
```

## Running Locally

Run the worker directly:

```
go run main.go
```

To specify the number of worker goroutines:

```
NUM_WORKERS=16 go run main.go
```

## Docker Compose

The worker is configured to run with Docker Compose:

```
docker compose up
```

All environment variables are configured in the docker-compose.yml file.

## Configuration Options

| Environment Variable | Default | Description |
|----------------------|---------|-------------|
| DB_HOST | localhost | Database host address |
| DB_PORT | 3306 | Database port |
| DB_USER | root | Database username |
| DB_PASSWORD | password | Database password |
| DB_NAME | pastebin | Database name |
| RABBITMQ_URL | amqp://guest:guest@localhost:5672/ | RabbitMQ connection URL |
| NUM_WORKERS | 10 | Number of worker goroutines to run |

## Worker Configuration

- RabbitMQ Queue: view_updates_queue
- RabbitMQ Exchange: view_updates (direct exchange)
- RabbitMQ Routing Key: view.update 