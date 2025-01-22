package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	// Prompt user for their name
	var name string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)

	// Listen for incoming messages from the chatroom
	go listenForMessages(natsConn)

	// Start reading user input and send messages
	for {
		var message string
		fmt.Print("You: ")
		fmt.Scanln(&message)

		// Send message to NATS server
		if message != "" {
			natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s: %s", name, message)))
		}
	}
}

func listenForMessages(natsConn *nats.Conn) {
	_, err := natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		// Print incoming messages to the console
		fmt.Printf("\nReceived message: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS channel:", err)
	}
}
