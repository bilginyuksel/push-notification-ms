package repository

import (
	"database/sql"
	"log"

	"github.com/bilginyuksel/push-notification/entity"
)

type TopicRepository interface {
	Save(t entity.Topic) error
	FindByAppIDAndName(appID, name string) (*entity.Topic, error)
}

type topicRepoImpl struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) TopicRepository {
	return &topicRepoImpl{db: db}
}

func (repo *topicRepoImpl) Save(t entity.Topic) error {
	query := "INSERT INTO C_APP_TOPIC (RECORD_ID, APP_ID, NAME, DESCRIPTION) VALUES (?, ?, ?, ?)"

	if res, err := repo.db.Exec(query,
		t.RecordID, t.AppID, t.Name, t.Description); err != nil {

		log.Printf("error occured while creating a topic, err: %v", err)

		return err
	} else {
		log.Printf("topic created, sql.result: %v", res)

		return nil
	}

}

func (repo *topicRepoImpl) FindByAppIDAndName(appID, name string) (*entity.Topic, error) {

	topic := &entity.Topic{}
	query := "SELECT RECORD_ID, APP_ID, NAME, DESCRIPTION FROM C_APP_TOPIC WHERE APP_ID=? AND NAME=?"

	if err := repo.db.QueryRow(query, appID, name).Scan(
		&topic.RecordID,
		&topic.AppID,
		&topic.Name,
		&topic.Description); err != nil {

		log.Printf("error found in db row scan operation, err: %v", err)
		return nil, err
	}

	return topic, nil
}
