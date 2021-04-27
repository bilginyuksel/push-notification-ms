package queue

import (
	"log"

	"github.com/Shopify/sarama"
)

type saramaProducer struct {
	sp *sarama.SyncProducer
}

func NewProducer(brokers []string) Producer {

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	sp, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Panicf("producer couldn't created, err: %v", err)
		return nil
	}

	return &saramaProducer{sp: &sp}
}

func (p *saramaProducer) Send(message string) error {
	return p.SendWithTopic(defaultTopic, message)
}

func (p *saramaProducer) SendWithTopic(topic, message string) error {
	partition, offset, err := (*p.sp).SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	})

	log.Printf("partition: %d, offset: %d", partition, offset)

	return err
}

func (p *saramaProducer) Close() error {
	if err := (*p.sp).Close(); err != nil {
		log.Printf("error occurred while closing producer connection, err: %v", err)
		return err
	}

	log.Println("producer closed successfully")
	return nil
}

func (p *saramaProducer) Ping() {

}
