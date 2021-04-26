package repository

import (
	"database/sql"

	"log"

	"github.com/bilginyuksel/push-notification/entity"
)

type appRepoImpl struct {
	db *sql.DB
}

type APPRepository interface {
	Save(a entity.Application) error
	IsExist(appID string) bool
}

func NewAPPRepository(db *sql.DB) APPRepository {
	return &appRepoImpl{db: db}
}

func (repo *appRepoImpl) Save(a entity.Application) error {
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

func (repo *appRepoImpl) IsExist(appID string) bool {
	var count int
	query := "SELECT COUNT(*) FROM C_APP WHERE RECORD_ID=?"

	if err := repo.db.QueryRow(query, appID).Scan(&count); err != nil {
		log.Printf("sql exception occurred, err: ", err)
	}

	return count != 0
}
