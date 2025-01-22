package natsclient

import (
	"log"

	"github.com/nats-io/nats.go"
)

func Connect() (*nats.Conn, error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
		return nil, err
	}
	return conn, nil
}
