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

	// Log when server successfully connects to NATS
	fmt.Println("Server connected to NATS.")

	// Subscribe to the 'chatroom' channel to listen for incoming messages
	_, err = natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		fmt.Printf("New message: %s\n", string(msg.Data))
	})

	if err != nil {
		log.Fatal("Error subscribing to NATS channel:", err)
	}

	// Log that the server is listening for messages
	fmt.Println("Server is listening for messages on 'chatroom'...")

	// Keep the server running to listen for messages
	select {}
}
