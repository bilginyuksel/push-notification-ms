package service

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/bilginyuksel/push-notification/entity"
)

type NotificationService interface {
	Push(notification entity.Notification)
}

type notificationServiceImpl struct {
	producer Producer
}

func NewNotificationService() NotificationService {
	brokers := []string{"localhost:9092"}

	return &notificationServiceImpl{
		producer: newSaramaProducer(brokers),
	}
}

func (service *notificationServiceImpl) Push(notification entity.Notification) {
}

type Producer interface {
	Send(topic, message string)
	Ping()
	Close() error
}

type producerImpl struct {
	sp *sarama.SyncProducer
}

func newSaramaProducer(brokers []string) Producer {

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	saramaProducer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Panicf("producer couldn't created, err: %v", err)
		return nil
	}

	return &producerImpl{sp: &saramaProducer}
}

func (p *producerImpl) Send(topic, message string) {
	if partition, offset, err := (*p.sp).SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}); err != nil {
		log.Printf("error occurred when sending message, err: %v", err)
	} else {
		log.Printf("message sent, partition: %d, offset: %d", partition, offset)
	}
}

func (p *producerImpl) Close() error {
	if err := (*p.sp).Close(); err != nil {
		log.Printf("error occurred while closing producer connection, err: %v", err)
		return err
	}

	log.Println("producer closed successfully")
	return nil
}

func (p *producerImpl) Ping() {

}
