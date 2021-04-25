package main

import (
	"database/sql"
	"log"
)

type clientRepoImpl struct {
	db *sql.DB
}

type ClientRepository interface {
	save(c Client) error
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepoImpl{db: db}
}

func (repo *clientRepoImpl) save(c Client) error {
	query := "INSERT INTO C_APP_USER (RECORD_ID, CLIENT_ID, APP_ID, STATUS, REGISTER_TIME, LAST_STATUS_CHANGE_TIME) VALUES (?, ?, ?, ?, ?, ?)"
	if res, err := repo.db.Exec(query,
		c.RecordID, c.ClientID, c.APPID, c.Status, c.RegisterTime, c.LastStatusChangeTime); err != nil {

		log.Printf("error occurred on client registering, err: %v", err)

		return err
	} else {
		log.Printf("client registered, sql result: %v", res)

		return nil
	}
}
