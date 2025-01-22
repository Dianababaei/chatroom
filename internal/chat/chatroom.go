package chat

import (
	"chatroom/internal/users"
	"fmt"

	"github.com/nats-io/nats.go"
)

type ChatRoom struct {
	natsConn *nats.Conn
}

func NewChatRoom(natsConn *nats.Conn) *ChatRoom {
	return &ChatRoom{natsConn: natsConn}
}

// Listen for messages and handle user interactions
func (c *ChatRoom) ListenForMessages(userManager *users.UserManager) {
	// Listen for new messages on NATS
	subscription, err := c.natsConn.Subscribe("chatroom", func(msg *nats.Msg) {
		fmt.Printf("New message: %s\n", string(msg.Data))
	})
	if err != nil {
		fmt.Println("Error subscribing to NATS:", err)
		return
	}
	defer subscription.Unsubscribe()

	// Broadcast user list on a request
	c.natsConn.Subscribe("users", func(msg *nats.Msg) {
		activeUsers := userManager.GetActiveUsers()
		c.natsConn.Publish(msg.Reply, []byte(fmt.Sprintf("Active Users: %v", activeUsers)))
	})
}

// Send a message to the chatroom via NATS
func (c *ChatRoom) SendMessage(message string) {
	c.natsConn.Publish("chatroom", []byte(message))
}
