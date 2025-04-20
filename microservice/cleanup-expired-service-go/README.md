# Cleanup Expired Pastes Service

This microservice is responsible for periodically deleting expired pastes from the database. It runs as a scheduled job at intervals defined in the environment variables.

## Features

- Automatically deletes pastes that have passed their expiration date
- Configurable cleanup interval via environment variables
- Runs as a standalone microservice within the Pastebin SOA architecture

## Configuration

The service can be configured using the following environment variables:

- `DB_HOST`: Database host (default: "db")
- `DB_PORT`: Database port (default: 3306)
- `DB_USER`: Database username
- `DB_PASS`: Database password
- `DB_NAME`: Database name (default: "pastebin")
- `CLEANUP_INTERVAL_MINUTES`: Interval in minutes between cleanup operations (default: 60)

## How It Works

1. The service connects to the database on startup
2. It immediately performs an initial cleanup of expired pastes
3. It then sets up a scheduler to run the cleanup operation at the configured interval
4. Each cleanup operation deletes all pastes where `expires_at` is not null and less than the current time

## Integration with Docker Compose

The service is integrated into the docker-compose.yml file and will start automatically with the rest of the application stack.

## Logs

The service logs the following information:
- Service startup
- Database connection status
- Number of pastes deleted during each cleanup operation
- Any errors encountered during operation
