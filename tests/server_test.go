package tests

import (
	"testing"
	"time"

	"chatroom/internal/users"

	"github.com/nats-io/nats.go"
)

func TestServer(t *testing.T) {
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		t.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	userManager := users.NewUserManager()

	natsConn.Subscribe("user.join", func(msg *nats.Msg) {
		userManager.AddUser(&users.User{Name: string(msg.Data)})
	})
	natsConn.Subscribe("user.leave", func(msg *nats.Msg) {
		userManager.RemoveUser(string(msg.Data))
	})

	natsConn.Publish("user.join", []byte("testuser1"))
	time.Sleep(100 * time.Millisecond)

	if len(userManager.GetActiveUsers()) != 1 {
		t.Errorf("Expected 1 active user, got %v", len(userManager.GetActiveUsers()))
	}

	natsConn.Publish("user.leave", []byte("testuser1"))
	time.Sleep(100 * time.Millisecond)

	if len(userManager.GetActiveUsers()) != 0 {
		t.Errorf("Expected 0 active users, got %v", len(userManager.GetActiveUsers()))
	}
}
