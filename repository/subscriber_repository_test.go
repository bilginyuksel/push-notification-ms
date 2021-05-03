package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/bilginyuksel/push-notification/entity"
)

func TestSubscriberRepository(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {
		tp := testPreperation{db: db}
		repo := NewSubscriberRepo(db)

		tp.createSampleApplication("test_app_id")
		tp.createSampleUser("test_user_id", "test_app_id")
		tp.createSampleTopic("test_topic_id", "test_app_id", "test_topic_name")

		subs := entity.Subscription{
			RecordID: "test_subs_id",
			AppID:    "test_app_id",
			UserID:   "test_user_id",
			TopicID:  "test_topic_id",
		}

		if err := repo.Save(subs); err != nil {
			t.Errorf("subscription creation failed, err: %v", err)
		}

		// Check the subscription from DB
		subsFromDb := entity.Subscription{}
		query := "SELECT RECORD_ID, APP_ID, USER_ID, TOPIC_ID FROM C_APP_TOPIC_USER WHERE RECORD_ID=?"
		if err := db.QueryRow(query, "test_subs_id").Scan(
			&subsFromDb.RecordID,
			&subsFromDb.AppID,
			&subsFromDb.UserID,
			&subsFromDb.TopicID,
		); err != nil {
			t.Errorf("subscription could not found, err: %v", err)
		}

		if !reflect.DeepEqual(subsFromDb, subs) {
			t.Errorf("subs and db are not equal, given: %v\nexpected: %v", subsFromDb, subs)
		}
	})
}
