# Pastebin Clone

A simple pastebin clone built with Node.js, Express, and MySQL.

## Environment Setup

The application uses Docker Compose for easy development and deployment. All environment variables are configured in the `.env` file and automatically used by both the application and Docker Compose.

### Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
NODE_ENV=development
PORT=3000

# Database Configuration
DB_HOST=mysql
DB_USER=root
DB_PASSWORD=pastebin123
DB_NAME=pastebin
```

## Running the Application

1. Make sure you have Docker and Docker Compose installed
2. Clone the repository
3. Copy `.env.example` to `.env` and adjust values if needed
4. Run the application:

```bash
docker-compose up
```

The application will be available at http://localhost:3000

## Development

- The application source code is in the `src` directory
- Database schema is in `src/database/schema.sql`
- The MySQL database is exposed on port 3307 for local development
- Source code changes will automatically reload the application
- Environment variables from `.env` are automatically used by both the application and Docker containers

## Cài đặt

1. Copy file môi trường:
```bash
cp .env.example .env
```

2. Build Docker images:
```bash
docker-compose build
```

3. Khởi chạy containers:
```bash
docker-compose up
```

Hoặc chạy ở chế độ detached (chạy ngầm):
```bash
docker-compose up -d
```

Để dừng các containers:
```bash
docker-compose down
``` 