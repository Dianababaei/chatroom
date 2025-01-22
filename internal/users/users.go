package users

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

// ChatRoomInterface defines the methods the User struct needs from a ChatRoom.
type ChatRoomInterface interface {
	SendMessage(message string)
}

// User represents a chatroom user.
type User struct {
	Name string
}

// UserManager manages the active users.
type UserManager struct {
	users map[string]*User
	mu    sync.Mutex
}

// Global variable to store user manager instance
var userManager *UserManager

// NewUserManager initializes the UserManager instance if not already done.
func NewUserManager() *UserManager {
	if userManager == nil {
		userManager = &UserManager{users: make(map[string]*User)}
	}
	return userManager
}

func (um *UserManager) AddUser(user *User) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.users[user.Name] = user
}

func (um *UserManager) RemoveUser(name string) {
	um.mu.Lock()
	defer um.mu.Unlock()
	delete(um.users, name)
}

func (um *UserManager) GetActiveUsers() []string {
	um.mu.Lock()
	defer um.mu.Unlock()

	var activeUsers []string
	for _, user := range um.users {
		activeUsers = append(activeUsers, user.Name)
	}
	return activeUsers
}

// NewUser creates a new user and asks for their name.
func NewUser() *User {
	var name string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)
	return &User{Name: name}
}

// ListenForInput handles user input and sends messages to the ChatRoomInterface.
func (u *User) ListenForInput(chatRoom ChatRoomInterface, natsConn *nats.Conn) {
	// Add user to the UserManager
	userManager.AddUser(u)

	// Remove user when they exit
	defer userManager.RemoveUser(u.Name)

	// Listen for the response to #users command
	go listenForUserListResponse(natsConn)

	// Main loop for user input
	for {
		var message string
		fmt.Print("You: ")
		fmt.Scanln(&message)

		if message == "" {
			// Skip empty messages
			continue
		}

		if message == "#users" {
			// Send request for active users
			natsConn.Publish("users", []byte("Requesting active users"))
		} else {
			// Send the message to chatroom
			chatRoom.SendMessage(fmt.Sprintf("%s: %s", u.Name, message))
		}
	}
}

// listenForUserListResponse listens for the server's response with the list of active users.
func listenForUserListResponse(natsConn *nats.Conn) {
	// Subscribe to the response channel from the server
	_, err := natsConn.Subscribe("users", func(msg *nats.Msg) {
		// Handle response from the server with active users list
		fmt.Printf("\nActive Users: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS users channel:", err)
	}
}
