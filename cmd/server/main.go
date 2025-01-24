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
	// Connect to the NATS server
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	// Initialize the UserManager
	userManager = users.NewUserManager()

	// Handle user join
	_, err = natsConn.Subscribe("user.join", func(msg *nats.Msg) {
		userName := string(msg.Data)
		userManager.AddUser(&users.User{Name: userName})
		log.Printf("User joined: %s", userName)
		natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s joined the chatroom.", userName)))
	})
	if err != nil {
		log.Fatal("Error subscribing to user.join channel:", err)
	}

	// Handle user leave
	_, err = natsConn.Subscribe("user.leave", func(msg *nats.Msg) {
		userName := string(msg.Data)
		userManager.RemoveUser(userName)
		log.Printf("User left: %s", userName)
		natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s left the chatroom.", userName)))
	})
	if err != nil {
		log.Fatal("Error subscribing to user.leave channel:", err)
	}

	// Respond to #users command
	_, err = natsConn.Subscribe("users", func(msg *nats.Msg) {
		activeUsers := userManager.GetActiveUsers()
		response := fmt.Sprintf("Active Users: %s", strings.Join(activeUsers, ", "))
		if msg.Reply != "" {
			natsConn.Publish(msg.Reply, []byte(response))
		}
	})
	if err != nil {
		log.Fatal("Error subscribing to users channel:", err)
	}

	log.Println("Server is running...")
	select {}
}
