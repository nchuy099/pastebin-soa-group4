# Pastebin Clone

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