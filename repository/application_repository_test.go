package repository

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/bilginyuksel/push-notification/entity"
)

func TestAppRepositorySave(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {

		repo := NewAPPRepository(db)

		application := entity.Application{
			RecordID:    "recordid",
			Name:        "Test_Application",
			Description: "Test_app_desc",
			Status:      "Active",
			CreateTime:  time.Now(),
			CancelTime:  nil,
		}

		if err := repo.Save(application); err != nil {
			t.Errorf("application save operation failed, err: %v", err)
		}

		appFromDB := entity.Application{}
		query := "SELECT RECORD_ID,NAME,DESCRIPTION,STATUS,CREATE_TIME,CANCEL_TIME FROM C_APP WHERE RECORD_ID=?"
		if err := db.QueryRow(query, application.RecordID).Scan(
			&appFromDB.RecordID,
			&appFromDB.Name,
			&appFromDB.Description,
			&appFromDB.Status,
			&appFromDB.CreateTime,
			&appFromDB.CancelTime,
		); err != nil {
			t.Errorf("scan operation failed, err: %v", err)
		}

		if reflect.DeepEqual(appFromDB, application) {
			t.Errorf("given and expected app are not matched")
		}

	})
}

func TestAppRepositoryDelete(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {
		repo := NewAPPRepository(db)

		application := entity.Application{
			RecordID:    "recordid",
			Name:        "Test_Application",
			Description: "Test_app_desc",
			Status:      "Active",
			CreateTime:  time.Now(),
			CancelTime:  nil,
		}

		if err := repo.Save(application); err != nil {
			t.Errorf("application save operation failed, err: %v", err)
		}

		// we believe this is created because we have another test for creation above
		err := repo.Delete("recordid")

		if err != nil {
			t.Errorf("an error occurred while deleting the record")
		}

		// control db
		query := "SELECT COUNT(*) FROM C_APP WHERE RECORD_ID=recordid"

		var count int
		if db.QueryRow(query).Scan(&count); err != nil || count != 0 {
			t.Errorf("there is a record in database with the same record id")
		}

	})
}
