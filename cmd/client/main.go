package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	// Prompt user for their name
	var name string
	fmt.Print("Enter your name: ")

	reader := bufio.NewReader(os.Stdin)
	name, _ = reader.ReadString('\n')
	name = name[:len(name)-1] // Remove newline character

	// Listen for incoming messages from the chatroom
	go listenForMessages(natsConn)

	// Handle system interrupt for graceful exit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start reading user input and send messages asynchronously
	go handleUserInput(natsConn, reader, name)

	// Wait for exit signal
	select {
	case <-signalChan:
		fmt.Println("Exiting gracefully...")
		return
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

	// Subscribe to the 'users' topic to listen for active user requests
	_, err = natsConn.Subscribe("users", func(msg *nats.Msg) {
		// Here, you would respond to the active users request.
		activeUsers := "Alice, Bob, Diana, Sadra"
		natsConn.Publish(msg.Reply, []byte(fmt.Sprintf("Active Users: %s", activeUsers)))
	})
	if err != nil {
		log.Fatal("Error subscribing to users topic:", err)
	}
}

func handleUserInput(natsConn *nats.Conn, reader *bufio.Reader, name string) {
	for {
		fmt.Print("You: ")
		message, _ := reader.ReadString('\n')
		message = message[:len(message)-1] // Remove newline character

		if message != "" {
			if message == "#users" {
				// Request the list of active users
				natsConn.Publish("users", []byte("Requesting active users"))
			} else {
				// Send message to the chatroom
				natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s: %s", name, message)))
			}
		}
	}
}
