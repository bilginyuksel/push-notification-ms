package main

import (
	"log"

	"github.com/Shopify/sarama"
)

var (
	brokers                 = []string{"localhost:9092"}
	topic                   = "sarama"
	partition               = 0
	activeProducer Producer = nil
)

type producerImpl struct {
	sp *sarama.SyncProducer
}

type Producer interface {
	SendMessage(topic, message string)
	Close() error
}

// GetProducer returns the active producer, if a producer is already created it will return that
// if there are no producers it will create the new one
func GetProducer() Producer {
	if activeProducer != nil {
		return activeProducer
	}

	activeProducer = newProducer()
	return activeProducer
}

func (p *producerImpl) SendMessage(topic, message string) {

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

func newProducer() Producer {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	saramaProducer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Panicf("producer couldn't created, err: %v", err)
		return nil
	}

	return &producerImpl{
		sp: &saramaProducer}

}
