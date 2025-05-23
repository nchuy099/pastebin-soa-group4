## Pastebin Scaling System
- A scalable Pastebin handling high concurrency with low p99 latency and minimal errors, verified by
extensive load testing.

- Implemented by SOA-Group4, members:
  - Nguyễn Chí Huy
  - Trần Quang Huy
  - Bùi Đức Huy
  - Nguyễn Đức Thành
### Technologies
- Go, ReactJs, Nginx, Redis, RabbitMQ, MySQL, Docker, Prometheus, Grafana, Loki
### Architecture
![image](https://github.com/user-attachments/assets/2febad2b-2f37-4a86-af12-2059d6ea6526)

### Load Testing
- Tool: Locust
- Machine: 6 vCPU, 6GB RAM
- Test Scenarios:

  ![Screenshot 2025-05-23 184502](https://github.com/user-attachments/assets/00dccad2-43b5-434f-87bc-3a87abf49fb5)
- Test Results

  ![Screenshot 2025-05-23 185123](https://github.com/user-attachments/assets/5d357181-8a94-40af-b507-a9866ac388b0)
