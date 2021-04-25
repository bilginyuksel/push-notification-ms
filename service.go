package main

import (
	"errors"
	"fmt"
)

func pushNotificationToClient(clientID string, notification NotificationRequest) {
	p := GetProducer()
	p.SendMessage(topic, fmt.Sprintf("cid:%s,%v", clientID, notification))
}

func pushNotificationToTopic(topic string, notification NotificationRequest) error {
	return errors.New("error")
}

func pushNotificationToApp(topic string, notification NotificationRequest) error {
	return errors.New("error")
}
