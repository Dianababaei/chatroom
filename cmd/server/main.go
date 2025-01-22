package main

import (
	"chatroom/internal/users"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

var userManager *users.UserManager

func main() {
	log.Println("Server is starting...")

	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	log.Println("Connected to NATS server")

	userManager = users.NewUserManager()

	// Subscribe to the 'chatroom' channel to listen for incoming messages
	_, err = natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		// Broadcast incoming messages and log them
		log.Printf("Received message on 'chatroom' channel: %s\n", string(msg.Data))
		fmt.Printf("New message: %s\n", string(msg.Data))
	})

	if err != nil {
		log.Fatal("Error subscribing to NATS 'chatroom' channel:", err)
	}

	log.Println("Subscribed to 'chatroom' channel for receiving messages")

	// Subscribe to the 'users' channel to send active user list when requested
	_, err = natsConn.Subscribe("users", func(msg *nats.Msg) {
		// Log when a request for active users is received
		log.Printf("Received request for active users: %s\n", msg.Data)

		// Fetch active users and send them back as a response
		activeUsers := userManager.GetActiveUsers()
		response := fmt.Sprintf("Active Users: %v", activeUsers)
		natsConn.Publish(msg.Reply, []byte(response))

		log.Printf("Sent active users response: %s\n", response)
	})

	if err != nil {
		log.Fatal("Error subscribing to NATS 'users' channel:", err)
	}

	log.Println("Subscribed to 'users' channel for responding with active users")

	log.Println("Server is running and waiting for messages...")
	select {}
}
