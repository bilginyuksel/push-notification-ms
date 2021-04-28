package repository

import (
	"database/sql"
	"log"

	"github.com/bilginyuksel/push-notification/entity"
)

type SubscriberRepository interface {
	Save(subscription entity.Subscription) error
	FindAllByAppIDAndTopicID(appID, topicID string) []*entity.Subscription
	FindAllByAppIDAndUserID(appID, userID string) []*entity.Subscription
	FindAllByUserID(userID string) []*entity.Subscription
}

type subscriberRepoImpl struct {
	db *sql.DB
}

func NewSubscriberRepo(db *sql.DB) SubscriberRepository {
	return &subscriberRepoImpl{db: db}
}

func (repo *subscriberRepoImpl) Save(subs entity.Subscription) error {
	query := "INSERT INTO C_APP_TOPIC_USER (RECORD_ID, APP_ID, TOPIC_ID, USER_ID) VALUES (?, ?, ?, ?)"

	if res, err := repo.db.Exec(query,
		subs.RecordID,
		subs.AppID,
		subs.TopicID,
		subs.UserID,
	); err != nil {
		log.Printf("error occurred while creating a subscription, err: %v", err)
		return err
	} else {
		log.Printf("subscription successfully saved, sql.res: %v", res)
		return nil
	}
}

func (repo *subscriberRepoImpl) FindAllByAppIDAndTopicID(appID, topicID string) []*entity.Subscription {
	query := "SELECT RECORD_ID, APP_ID, TOPIC_ID, USER_ID FROM C_APP_TOPIC_USER WHERE APP_ID=? AND TOPIC_ID=?"

	rows, err := repo.db.Query(query, appID, topicID)

	return repo.getSubscriptions(rows, err)
}

func (repo *subscriberRepoImpl) FindAllByAppIDAndUserID(appID, userID string) []*entity.Subscription {
	query := "SELECT RECORD_ID, APP_ID, TOPIC_ID, USER_ID FROM C_APP_TOPIC_USER WHERE APP_ID=? AND USER_ID=?"

	rows, err := repo.db.Query(query, appID, userID)

	return repo.getSubscriptions(rows, err)
}
func (repo *subscriberRepoImpl) FindAllByUserID(userID string) []*entity.Subscription {
	query := "SELECT RECORD_ID, APP_ID, TOPIC_ID, USER_ID FROM C_APP_TOPIC_USER WHERE USER_ID=?"

	rows, err := repo.db.Query(query, userID)

	return repo.getSubscriptions(rows, err)
}

func (repo *subscriberRepoImpl) getSubscriptions(rows *sql.Rows, err error) []*entity.Subscription {
	subscriptions := []*entity.Subscription{}

	if err != nil {
		log.Printf("error occurred while querying subscriptions, err: %v", err)
		return subscriptions
	}

	for rows.Next() {
		subs := &entity.Subscription{}
		if err := rows.Scan(
			&subs.RecordID,
			&subs.AppID,
			&subs.TopicID,
			&subs.UserID,
		); err != nil {
			log.Printf("error ocurred while scanning a subscription row, err: %v", err)
		} else {
			subscriptions = append(subscriptions, subs)
		}
	}

	return subscriptions
}
