package main

import (
	"chatroom/internal/users"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

var userManager *users.UserManager

func main() {
	// Connect to NATS server
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	// Initialize user manager
	userManager = users.NewUserManager()

	// Subscribe to the 'chatroom' channel for receiving messages
	_, err = natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		log.Printf("Received message on 'chatroom' channel: %s", string(msg.Data))
		fmt.Printf("\n%s\n", string(msg.Data))
	})

	if err != nil {
		log.Fatal("Error subscribing to NATS channel:", err)
	}

	// Subscribe to the 'users' channel to handle active user requests
	_, err = natsConn.Subscribe("users", func(msg *nats.Msg) {
		// Get the active user list and respond to the requesting client
		activeUsers := userManager.GetActiveUsers()
		response := fmt.Sprintf("Active Users: %v", activeUsers)
		if msg.Reply != "" {
			natsConn.Publish(msg.Reply, []byte(response))
		}
	})

	if err != nil {
		log.Fatal("Error subscribing to NATS users channel:", err)
	}

	log.Println("Server is running and waiting for messages...")
	select {}
}
