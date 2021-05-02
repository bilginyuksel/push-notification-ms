package repository

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/bilginyuksel/push-notification/entity"
)

const (
	cAppInsertQuery = "INSERT INTO C_APP (RECORD_ID, NAME, DESCRIPTION, STATUS, CREATE_TIME, CANCEL_TIME) VALUES ('recordid', 'test_app', 'test_desc', 'Active', '2021-05-01 08:04:15', NULL)"
)

func TestAppRepositorySave(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {

		repo := NewAPPRepository(db)

		application := entity.Application{
			RecordID:    "recordid",
			Name:        "Test_Application",
			Description: "Test_app_desc",
			Status:      "Active",
			CreateTime:  time.Now().UTC().Round(time.Second),
			CancelTime:  nil,
		}

		if err := repo.Save(application); err != nil {
			t.Errorf("application save operation failed, err: %v", err)
		}

		// Query the new inserted Application to test
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

		if !reflect.DeepEqual(appFromDB, application) {
			t.Errorf("given and expected app are not matched, expected: %v\n fromDb: %v", application, appFromDB)
		}

	})
}

func TestAppRepositoryDelete(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {
		repo := NewAPPRepository(db)

		_, err := db.Exec(cAppInsertQuery)
		if err != nil {
			t.Errorf("insert query error at testing preperation phase, err: %v", err)
		}

		err = repo.Delete("recordid")

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

func TestAppRepositoryIsExist(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {
		repo := NewAPPRepository(db)

		_, err := db.Exec(cAppInsertQuery)
		if err != nil {
			t.Errorf("insert query error at testing preperation phase, err: %v", err)
		}

		if exists := repo.IsExist("recordid"); !exists {
			t.Errorf("record inserted, but the function returned not exist")
		}
		if exists := repo.IsExist("nonoo"); exists {
			t.Errorf("record not inserted, but the function returned exist")
		}
	})
}
