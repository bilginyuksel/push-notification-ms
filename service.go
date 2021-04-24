package main

import (
	"errors"
	"fmt"
)

func pushNotificationToClient(clientID string, notification PushNotificationReq) error {
	p := GetProducer()
	p.SendMessage(topic, fmt.Sprintf("cid:%s,%v", clientID, notification))
	return nil
}

func pushNotificationToTopic(topic string, notification PushNotificationReq) error {
	return errors.New("error")
}

func pushNotificationToApp(topic string, notification PushNotificationReq) error {
	return errors.New("error")
}
