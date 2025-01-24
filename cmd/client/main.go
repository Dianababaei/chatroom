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
	// Connect to the NATS server
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	// Prompt user for their name
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Notify the server of user join
	natsConn.Publish("user.join", []byte(name))
	defer natsConn.Publish("user.leave", []byte(name)) // Notify server when user leaves

	// Create a reply subject for the `#users` command
	replySubject := fmt.Sprintf("users.%s", name)
	go listenForUserList(natsConn, replySubject)

	// Subscribe to chatroom messages
	go listenForMessages(natsConn, name)

	// Display simple CLI instructions
	fmt.Println("\nWelcome to the Chatroom!")
	fmt.Println("Commands:")
	fmt.Println("  Type your message and press Enter to send.")
	fmt.Println("  Type '#users' to see the list of active users.")
	fmt.Println("  Type '#exit' to leave the chatroom.")
	fmt.Println()

	// Declare the signal channel to handle system interrupts
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Main CLI loop
	for {
		fmt.Print("You: ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		switch command {
		case "#users":
			// Request the list of active users
			natsConn.PublishRequest("users", replySubject, nil)
			// Wait for the response to print active users
		case "#exit":
			// Exit the chatroom
			fmt.Println("\nExiting the chatroom...")
			return
		default:
			// Send a chat message
			if command != "" {
				natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s: %s", name, command)))
			}
		}

		// Handle interrupt signal for exiting
		select {
		case <-signalChan:
			fmt.Println("\nExiting the chatroom...")
			return
		default:
		}
	}
}

// listenForMessages subscribes to the chatroom topic and prints incoming messages.
func listenForMessages(natsConn *nats.Conn, clientName string) {
	_, err := natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		message := string(msg.Data)
		if !strings.HasPrefix(message, clientName+":") {
			// Skip messages sent by this client
			fmt.Printf("\r%s\nYou: ", message)
		}
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS chatroom channel:", err)
	}
}

// listenForUserList subscribes to the `#users` reply subject only once.
func listenForUserList(natsConn *nats.Conn, replySubject string) {
	// Make sure that we are not repeatedly printing "Active Users"
	_, err := natsConn.Subscribe(replySubject, func(msg *nats.Msg) {
		fmt.Printf("\r %s\nYou: ", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to user list reply channel:", err)
	}
}
