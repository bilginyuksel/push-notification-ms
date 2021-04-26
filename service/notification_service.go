package service

import "github.com/bilginyuksel/push-notification/entity"

type NotificationService interface {
	Push(notification entity.Notification)
}

type notificationServiceImpl struct {
}

func NewNotificationService() NotificationService {
	return &notificationServiceImpl{}
}

func (service *notificationServiceImpl) Push(notification entity.Notification) {

}
