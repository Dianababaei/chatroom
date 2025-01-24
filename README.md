```markdown
 Chatroom Application

A simple real-time chatroom built in Go using NATS as the message broker. It supports public chat, private messages, and user management.

Key Features:
- Real-time messaging.
- Public chatroom and private messaging.
- View active users.
- User join/leave notifications.

 How to Run:

 Prerequisites:
- Docker and Docker Compose for containerization.
- Go to run the client manually.

 Setup:
1. Build containers:
    ```bash
    docker-compose build
    ```

2. Start services:
    ```bash
    docker-compose up
    ```

3. Run the client manually:
    ```bash
    go run cmd/client/main.go
    ```

### Stop services:
```bash
docker-compose down
```

## Technologies Used:
- **Go**, **NATS**, **Docker**, **Docker Compose**.

