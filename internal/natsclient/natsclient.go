package natsclient

import (
	"log"

	"github.com/nats-io/nats.go"
)

func Connect() (*nats.Conn, error) {
	// Connect to NATS server (default)
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
		return nil, err
	}
	return conn, nil
}
