package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestPrivateMessage(t *testing.T) {
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		t.Fatal("Error connecting to NATS server:", err)
	}
	defer natsConn.Close()

	sender := "testuser1"
	receiver := "testuser2"
	message := "Hello, testuser2!"
	privateSubject := fmt.Sprintf("private.%s", receiver)

	var wg sync.WaitGroup
	wg.Add(1)

	_, err = natsConn.Subscribe(privateSubject, func(msg *nats.Msg) {
		if string(msg.Data) != fmt.Sprintf("[Private] %s: %s", sender, message) {
			t.Errorf("Expected message to be '[Private] %s: %s', got: %s", sender, message, string(msg.Data))
		} else {
			t.Logf("Private message received: %s", string(msg.Data))
		}
		wg.Done()
	})
	if err != nil {
		t.Fatal("Error subscribing to private message channel:", err)
	}

	natsConn.Publish(privateSubject, []byte(fmt.Sprintf("[Private] %s: %s", sender, message)))

	done := make(chan bool, 1)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		t.Log("Private message received successfully.")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout: Private message not received")
	}
}
