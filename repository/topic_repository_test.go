package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/bilginyuksel/push-notification/entity"
)

func TestTopicRepositorySave(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {

		tp := testPreperation{db}
		repo := NewTopicRepository(db)

		// create app because it is foreign key that bounds these together
		tp.createSampleApplication("test_app_id")

		topic := entity.Topic{
			RecordID:    "test_id",
			AppID:       "test_app_id",
			Name:        "test_name",
			Description: "test_desc",
		}

		if err := repo.Save(topic); err != nil {
			t.Errorf("topic could not created, err: %v", err)
		}

		// Test if topic is updated in the database
		query := "SELECT RECORD_ID, APP_ID, NAME, DESCRIPTION FROM C_APP_TOPIC WHERE RECORD_ID=?"
		topicFromDb := entity.Topic{}
		if err := db.QueryRow(query, "test_id").Scan(
			&topicFromDb.RecordID,
			&topicFromDb.AppID,
			&topicFromDb.Name,
			&topicFromDb.Description,
		); err != nil {
			t.Errorf("getting test data from db failed, err: %v", err)
		}

		if !reflect.DeepEqual(topic, topicFromDb) {
			t.Errorf("expected topic and the topic from db is not equal, expected: %v\nfromDb: %v", topic, topicFromDb)
		}
	})
}

func TestTopicRepositoryFindByAppIDAndName(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {
		tp := testPreperation{db: db}
		repo := NewTopicRepository(db)

		// create app because it is foreign key that bounds these together
		tp.createSampleApplication("test_app_id")
		tp.createSampleTopic("test_topic_id", "test_app_id", "test_topic_name")

		topic, err := repo.FindByAppIDAndName("test_app_id", "test_topic_name")

		if err != nil {
			t.Errorf("topic could not find, err: %v", err)
		}

		expected := entity.Topic{
			RecordID:    "test_topic_id",
			AppID:       "test_app_id",
			Name:        "test_topic_name",
			Description: "test_topic_desc",
		}

		if !reflect.DeepEqual(*topic, expected) {
			t.Errorf("expected and given topics are not same, given: %v\n expected: %v", topic, expected)
		}
	})
}
