package messaging

type Config struct {
	User         string
	Password     string
	Host         string
}

type MessagingQueue interface {
	Publish(subject string, message []byte) error
}