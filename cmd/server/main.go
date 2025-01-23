package main

import (
	"chatroom/internal/users"
	"fmt"
	"log"
	"strings"

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

	// Initialize UserManager
	userManager = users.NewUserManager()

	// Subscribe to 'chatroom' to handle messages
	_, err = natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		log.Printf("Received message on 'chatroom' channel: %s", string(msg.Data))
		fmt.Printf("\n%s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS chatroom channel:", err)
	}

	// Subscribe to 'users' to respond with the active user list
	_, err = natsConn.Subscribe("users", func(msg *nats.Msg) {
		activeUsers := userManager.GetActiveUsers()
		response := fmt.Sprintf("Active Users: %s", strings.Join(activeUsers, ", "))
		if msg.Reply != "" {
			natsConn.Publish(msg.Reply, []byte(response))
		}
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS users channel:", err)
	}

	// Handle user join
	_, err = natsConn.Subscribe("user.join", func(msg *nats.Msg) {
		userManager.AddUser(&users.User{Name: string(msg.Data)})
		log.Printf("User joined: %s", string(msg.Data))
		natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s joined the chatroom.", string(msg.Data))))
	})
	if err != nil {
		log.Fatal("Error subscribing to user.join channel:", err)
	}

	// Handle user leave
	_, err = natsConn.Subscribe("user.leave", func(msg *nats.Msg) {
		userManager.RemoveUser(string(msg.Data))
		log.Printf("User left: %s", string(msg.Data))
		natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s left the chatroom.", string(msg.Data))))
	})
	if err != nil {
		log.Fatal("Error subscribing to user.leave channel:", err)
	}

	log.Println("Server is running and waiting for messages...")
	select {}
}
