package queue

var (
	defaultTopic   = "sarama"
	defaultBrokers = []string{"localhost:9092"}
)

type Producer interface {
	// Send send a message using default topic
	Send(message producerMessage) error

	// SendMessages send multiple messages
	SendMessages(message []producerMessage) error

	// Close producer connection
	Close() error
}

type producerMessage struct {
	topic   string
	message string
}

// NewMessage use this method to produce a producerMessage via topic
func NewMessageWithTopic(topic, message string) producerMessage {
	return producerMessage{
		topic:   topic,
		message: message,
	}
}

// NewMessage use this method to produce a producer message
func NewMessage(message string) producerMessage {
	return producerMessage{
		topic:   defaultTopic,
		message: message,
	}
}
