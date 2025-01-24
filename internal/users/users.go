package users

import (
	"sync"
)

// User represents a chatroom user.
type User struct {
	Name string
}

// UserManager manages the active users.
type UserManager struct {
	users map[string]*User
	mu    sync.Mutex
}

// NewUserManager creates a new instance of UserManager.
func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*User),
	}
}

// AddUser adds a user to the active list.
func (um *UserManager) AddUser(user *User) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.users[user.Name] = user
}

// RemoveUser removes a user from the active list.
func (um *UserManager) RemoveUser(name string) {
	um.mu.Lock()
	defer um.mu.Unlock()
	delete(um.users, name)
}

// GetActiveUsers returns a list of active users' names.
func (um *UserManager) GetActiveUsers() []string {
	um.mu.Lock()
	defer um.mu.Unlock()

	activeUsers := make([]string, 0, len(um.users))
	for _, user := range um.users {
		activeUsers = append(activeUsers, user.Name)
	}
	return activeUsers
}
