package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	reader := bufio.NewReader(os.Stdin)
	var name string
	for {
		fmt.Print("Enter your name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		// Send user.join request with reply subject
		replySubject := fmt.Sprintf("join.%s", name)

		// Subscribe to the reply subject for server response
		sub, err := natsConn.SubscribeSync(replySubject)
		if err != nil {
			log.Fatal("Error subscribing to join reply channel:", err)
		}

		// Publish the join request
		natsConn.PublishRequest("user.join", replySubject, []byte(name))

		// Wait for server response
		msg, err := sub.NextMsg(2 * time.Second)
		sub.Unsubscribe()

		if err != nil {
			fmt.Println("Server did not respond. Please try again.")
			continue
		}

		response := string(msg.Data)
		if response == "This name already exists." {
			fmt.Println(response)
			continue
		}

		break
	}

	go listenForMessages(natsConn, name)
	go listenForPrivateMessages(natsConn, name)

	// Create a persistent subscription for user list responses
	replySubject := fmt.Sprintf("users.%s", name)
	go listenForUserList(natsConn, replySubject)

	time.Sleep(90 * time.Millisecond)

	natsConn.Publish("user.join", []byte(name))
	defer natsConn.Publish("user.leave", []byte(name))

	fmt.Println("\nWelcome to the Chatroom!")
	fmt.Println("Commands:")
	fmt.Println("  Type your message and press Enter to send.")
	fmt.Println("  Type '#users' to see the list of active users.")
	fmt.Println("  Type '#msg <username> <message>' to send a private message.")
	fmt.Println("  Type '#exit' to leave the chatroom.")
	fmt.Println()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		fmt.Print("You: ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		switch {
		case command == "#users":
			natsConn.PublishRequest("users", replySubject, nil)

		case command == "#exit":
			fmt.Println("\nExiting the chatroom...")
			return

		default:
			if strings.HasPrefix(command, "#msg") {
				parts := strings.SplitN(command, " ", 3)
				if len(parts) < 3 {
					fmt.Println("Usage: #msg <username> <message>")
					continue
				}
				targetUser := parts[1]
				message := parts[2]
				privateSubject := fmt.Sprintf("private.%s", targetUser)
				natsConn.Publish(privateSubject, []byte(fmt.Sprintf("[Private] %s: %s", name, message)))
			} else {
				if command != "" {
					natsConn.Publish("chatroom", []byte(fmt.Sprintf("%s: %s", name, command)))
				}
			}
		}

		select {
		case <-signalChan:
			fmt.Println("\nExiting the chatroom...")
			return
		default:
		}
	}
}

func listenForMessages(natsConn *nats.Conn, clientName string) {
	_, err := natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		message := string(msg.Data)
		if !strings.HasPrefix(message, clientName+":") {
			fmt.Printf("\r%s\nYou: ", message)
		}
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS chatroom channel:", err)
	}
}

func listenForUserList(natsConn *nats.Conn, replySubject string) {
	// Persistent subscription for user list responses
	_, err := natsConn.Subscribe(replySubject, func(msg *nats.Msg) {
		fmt.Printf("\r%s\nYou: ", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to user list reply channel:", err)
	}
}

func listenForPrivateMessages(natsConn *nats.Conn, username string) {
	privateSubject := fmt.Sprintf("private.%s", username)
	_, err := natsConn.Subscribe(privateSubject, func(msg *nats.Msg) {
		fmt.Printf("\r%s\nYou: ", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to private message channel:", err)
	}
}
