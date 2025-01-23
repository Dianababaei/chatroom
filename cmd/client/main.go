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
	name = strings.TrimSpace(name)

	// Notify server that the user has joined
	natsConn.Publish("user.join", []byte(name))
	defer natsConn.Publish("user.leave", []byte(name)) // Notify server when the user leaves

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
			// Request the list of active users
			replySubject := fmt.Sprintf("users.%s", name)
			go listenForUserList(natsConn, replySubject)
			natsConn.PublishRequest("users", replySubject, []byte(""))
		} else {
			// Broadcast message to the chatroom
			natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s: %s", name, message)))
		}

		// Handle exit signal
		select {
		case <-signalChan:
			fmt.Println("\nExiting chatroom...")
			return
		default:
		}
	}
}

func listenForMessages(natsConn *nats.Conn, clientName string) {
	_, err := natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		message := string(msg.Data)
		if strings.HasPrefix(message, clientName+":") {
			// Skip messages sent by this client
			return
		}

		// Clear the current line, print the new message, and re-display "You: " for user input
		fmt.Printf("\r%s\nYou: ", message)
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS chatroom channel:", err)
	}
}

func listenForUserList(natsConn *nats.Conn, replySubject string) {
	_, err := natsConn.Subscribe(replySubject, func(msg *nats.Msg) {
		// Clear the current line, print the list of active users, and re-display "You: "
		fmt.Printf("\r%s\nYou: ", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to user list reply channel:", err)
	}
}
