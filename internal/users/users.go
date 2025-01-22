package users

import (
	"fmt"
	"sync"
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

func NewUserManager() *UserManager {
	return &UserManager{users: make(map[string]*User)}
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

func NewUser() *User {
	var name string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)
	return &User{Name: name}
}

// ListenForInput handles user input and sends messages to the ChatRoomInterface.
func (u *User) ListenForInput(chatRoom ChatRoomInterface) {
	for {
		var message string
		fmt.Print("You: ")
		fmt.Scanln(&message)
		if message == "#users" {
			fmt.Println("Command '#users' is not implemented yet in this example.")
		} else {
			chatRoom.SendMessage(fmt.Sprintf("%s: %s", u.Name, message))
		}
	}
}
