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
	// Save an application to db
	Save(a entity.Application) error

	// Query get list of applications
	GetAll(userID string) ([]*entity.Application, error)

	// Delete an application with the given id from db
	Delete(userID, recordID string) error

	// IsExist return true if application found with the given id
	IsExist(userID, recordID string) bool
}

func NewAPPRepository(db *sql.DB) APPRepository {
	return &appRepoImpl{db: db}
}

func (repo *appRepoImpl) Save(a entity.Application) error {
	query := "INSERT INTO C_APP (USER_ID, RECORD_ID, NAME, DESCRIPTION, STATUS, CREATE_TIME, CANCEL_TIME) VALUES (?, ?, ?, ?, ?, ?, ?)"
	if res, err := repo.db.Exec(query,
		a.UserID, a.RecordID, a.Name, a.Description, a.Status, a.CreateTime, a.CancelTime); err != nil {

		log.Printf("error occured on application creation, err: %v", err)

		return err
	} else {
		log.Printf("application created, sql.result: %v", res)

		return nil
	}
}

func (repo *appRepoImpl) IsExist(userID, recordID string) bool {
	var count int
	query := "SELECT COUNT(*) FROM C_APP WHERE USER_ID=? AND RECORD_ID=?"

	if err := repo.db.QueryRow(query, userID, recordID).Scan(&count); err != nil {
		log.Printf("sql exception occurred, err: %v", err)
	}

	return count != 0
}

func (repo *appRepoImpl) Delete(userID, recordID string) error {
	query := "DELETE FROM C_APP WHERE USER_ID=? AND RECORD_ID=?"

	if res, err := repo.db.Exec(query, userID, recordID); err != nil {
		log.Printf("error occurred while deleting application, err: %v", err)

		return err
	} else {
		log.Printf("application deleted, sql.result: %v", res)

		return nil
	}
}

func (repo *appRepoImpl) GetAll(userID string) ([]*entity.Application, error) {
	query := "SELECT USER_ID, RECORD_ID, NAME, DESCRIPTION, STATUS, CREATE_TIME, CANCEL_TIME FROM C_APP WHERE USER_ID=?"

	appList := []*entity.Application{}

	rows, err := repo.db.Query(query, userID)

	if err != nil {
		log.Printf("error occurred while getting app, err: %v", err)
		return nil, err
	}

	for rows.Next() {
		app := &entity.Application{}
		if err := rows.Scan(
			&app.UserID,
			&app.RecordID,
			&app.Name,
			&app.Description,
			&app.Status,
			&app.CreateTime,
			&app.CancelTime,
		); err != nil {
			log.Printf("error occurred on scanning query result")
		} else {
			appList = append(appList, app)
		}
	}

	return appList, nil
}
