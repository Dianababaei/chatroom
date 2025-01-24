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
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	userManager = users.NewUserManager()

	_, err = natsConn.Subscribe("user.join", func(msg *nats.Msg) {
		userName := string(msg.Data)

		if userManager.IsUserExists(userName) {
			if msg.Reply != "" {
				natsConn.Publish(msg.Reply, []byte("This name already exists."))
			}
			return
		}

		userManager.AddUser(&users.User{Name: userName})
		log.Printf("User joined: %s", userName)

		if msg.Reply != "" {
			natsConn.Publish(msg.Reply, []byte("Welcome to the chatroom!"))
		}

		natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s joined the chatroom.", userName)))
	})
	if err != nil {
		log.Fatal("Error subscribing to user.join channel:", err)
	}

	_, err = natsConn.Subscribe("user.leave", func(msg *nats.Msg) {
		userName := string(msg.Data)
		userManager.RemoveUser(userName)
		log.Printf("User left: %s", userName)
		natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s left the chatroom.", userName)))
	})
	if err != nil {
		log.Fatal("Error subscribing to user.leave channel:", err)
	}

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
