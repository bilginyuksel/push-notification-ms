package service

import (
	"encoding/json"
	"log"

	"github.com/bilginyuksel/push-notification/entity"
	"github.com/bilginyuksel/push-notification/queue"
)

type NotificationService interface {
	// PushToTopic(topic string, notification entity.Notification)
	// PushToApp(appID string, notification entity.Notification)
	Push(notification entity.Notification) error
}

type notificationServiceImpl struct {
	producer queue.Producer
}

func NewNotificationService() NotificationService {
	brokers := []string{"localhost:9092"}

	return &notificationServiceImpl{
		producer: queue.NewProducer(brokers),
	}
}

func (service *notificationServiceImpl) Push(notification entity.Notification) error {
	bytes, err := json.Marshal(notification)

	if err != nil {
		log.Printf("an error occurred while marshaling the notification, err: %v", err)
		return err
	}

	message := queue.NewMessage(string(bytes))

	return service.producer.Send(message)
}
