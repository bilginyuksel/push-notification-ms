package repository

import (
	"database/sql"

	"github.com/bilginyuksel/push-notification/entity"
)

type TopicRepository interface {
	Save(t entity.Topic)
}

type topicRepoImpl struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) TopicRepository {
	return &topicRepoImpl{db: db}
}

func (repo *topicRepoImpl) Save(t entity.Topic) {

}
