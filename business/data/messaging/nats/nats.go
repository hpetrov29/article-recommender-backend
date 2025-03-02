package nats

import (
	"errors"
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/business/data/messaging"
	"github.com/nats-io/nats.go"
)

// NATSClient is a concrete implementation of MessageQueue.
type NATSClient struct {
    conn *nats.Conn
}

// NewNATSClient initializes a new NATSClient instance.
func NewNATSClient(cfg messaging.Config) (*NATSClient, error) {
	if cfg.Host == "" {
		cfg.Host = "nats-container:4222"
	}
	URI := fmt.Sprintf("nats://%s:%s@%s", cfg.User, cfg.Password, cfg.Host)
    nc, err := nats.Connect(URI)
    if err != nil {
        return nil, err
    }

    return &NATSClient{conn: nc}, nil
}

// Publish sends a message to a NATS subject.
func (n *NATSClient) Publish(subject string, message []byte) error {
    return n.conn.Publish(subject, message)
}

// HealthCheck verifies if the NATS connection is active.
func (n *NATSClient) StatusCheck() error {
	var status nats.Status
	for i := 0; i < 5; i++ {
		status = n.conn.Status()
		if (status == nats.CONNECTED) {
			return nil
		}
		time.Sleep(5*time.Second)
	}
    return errors.New(status.String())
}

// Close shuts down the NATS connection.
func (n *NATSClient) Close() {
    n.conn.Close()
}