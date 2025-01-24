package users

import (
	"sync"
)

type User struct {
	Name string
}

type UserManager struct {
	users map[string]*User
	mu    sync.Mutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*User),
	}
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

	activeUsers := make([]string, 0, len(um.users))
	for _, user := range um.users {
		activeUsers = append(activeUsers, user.Name)
	}
	return activeUsers
}
