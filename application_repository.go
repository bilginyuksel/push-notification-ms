package main

import (
	"database/sql"
	"log"
)

type appRepoImpl struct {
	db *sql.DB
}

type APPRepository interface {
	save(a Application) error
}

func NewAPPRepository(db *sql.DB) APPRepository {
	return &appRepoImpl{db: db}
}

func (repo *appRepoImpl) save(a Application) error {
	query := "INSERT INTO C_APP (RECORD_ID, NAME, DESCRIPTION, STATUS, CREATE_TIME, CANCEL_TIME) VALUES (?, ?, ?, ?, ?, ?)"
	if res, err := repo.db.Exec(query,
		a.RecordID, a.Name, a.Description, a.Status, a.CreateTime, a.CancelTime); err != nil {

		log.Printf("error occured on application creation, err: %v", err)

		return err
	} else {
		log.Printf("application created, sql.result: %v", res)

		return nil
	}
}
