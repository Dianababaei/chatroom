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
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	natsConn.Publish("user.join", []byte(name))
	defer natsConn.Publish("user.leave", []byte(name))

	go listenForPrivateMessages(natsConn, name)

	go listenForMessages(natsConn, name)

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
			replySubject := fmt.Sprintf("users.%s", name)
			go listenForUserList(natsConn, replySubject)
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
