package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name) // Remove extra spaces or newlines

	// Subscribe to a unique reply channel for this client
	replySubject := fmt.Sprintf("users.%s", name)
	go listenForUserList(natsConn, replySubject)

	// Listen for incoming chat messages
	go listenForMessages(natsConn, name)

	// Handle system interrupt for graceful exit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Main loop: Read user input and send messages
	for {
		fmt.Print("You: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}

		if message == "#users" {
			// Send request for active users
			natsConn.PublishRequest("users", replySubject, []byte("Requesting active users"))
		} else {
			// Broadcast message to the chatroom
			natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s: %s", name, message)))
		}

		// Handle exit signal (Ctrl+C)
		select {
		case <-signalChan:
			fmt.Println("\nExiting chatroom...")
			return
		default:
			// Keep running
		}
	}
}

func listenForMessages(natsConn *nats.Conn, clientName string) {
	_, err := natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		parts := strings.SplitN(string(msg.Data), ":", 2)
		if len(parts) < 2 {
			return
		}

		sender := parts[0]
		content := parts[1]

		if sender == clientName {
			return
		}

		// Clear the prompt, show the message, and reprint the prompt
		fmt.Printf("\r%s: %s\nYou: ", sender, content)
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS chatroom channel:", err)
	}
}

func listenForUserList(natsConn *nats.Conn, replySubject string) {
	_, err := natsConn.Subscribe(replySubject, func(msg *nats.Msg) {
		// Print the list of active users
		fmt.Printf("\n%s\nYou: ", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to user list reply channel:", err)
	}
}
