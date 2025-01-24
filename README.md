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

 Stop services:
```bash
docker-compose down
```

 Technologies Used:
- **Go**, **NATS**, **Docker**, **Docker Compose**.

- ............................

- Here's the corrected **README.md** format that you can directly use in GitHub:

```markdown
# Chatroom Application

A simple real-time chatroom built in **Go** using **NATS** as the message broker. It supports public chat, private messages, and user management.

## Key Features
- Real-time messaging.
- Public chatroom and private messaging.
- View active users.
- User join/leave notifications.

## How to Run

### Prerequisites
- **Docker** and **Docker Compose** for containerization.
- **Go** to run the client manually.

### Setup

1. **Build containers**:
    ```bash
    docker-compose build
    ```

2. **Start services**:
    ```bash
    docker-compose up
    ```

3. **Run the client manually**:
    ```bash
    go run cmd/client/main.go
    ```

4. **Interact with the chatroom**:
    - Type your message to send.
    - Use `#users` to see active users.
    - Use `#msg <username> <message>` to send private messages.
    - Type `#exit` to leave.

### Stop services
```bash
docker-compose down
```

## Known Issues
- **Manual client**: Needs to be run with `go run cmd/client/main.go`.
- **No authentication**.
- **No persistent message history**.
- **Limited scalability**.

## Technologies Used
- **Go**, **NATS**, **Docker**, **Docker Compose**.

## Conclusion
A basic real-time chatroom app with potential for improvement and expansion.
```

This version should now display correctly on GitHub. Just copy and paste it into your `README.md` file in your repository.

