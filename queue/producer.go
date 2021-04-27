package queue

var (
	defaultTopic   = "sarama"
	defaultBrokers = []string{"localhost:9092"}
)

type Producer interface {
	// SendWithTopic send a message using provided topic
	SendWithTopic(topic, message string) error

	// Send send a message using default topic
	Send(message string) error

	// Ping test connection
	Ping()

	// Close producer connection
	Close() error
}
