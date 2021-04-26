package service

import (
	"github.com/bilginyuksel/push-notification/entity"
	"github.com/bilginyuksel/push-notification/repository"
	"github.com/bilginyuksel/push-notification/request"
)

type topicServiceImpl struct {
	repo repository.TopicRepository
}

type TopicService interface {
	CreateNewTopic(req request.CreateTopicRequest) (*entity.Topic, error)
}

func NewTopicService(topicRepo repository.TopicRepository) TopicService {
	return &topicServiceImpl{
		repo: topicRepo,
	}
}

func (service topicServiceImpl) CreateNewTopic(req request.CreateTopicRequest) (*entity.Topic, error) {
	return nil, nil
}
