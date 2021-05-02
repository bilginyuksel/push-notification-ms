package repository

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/bilginyuksel/push-notification/entity"
)

func TestClientRepositorySave(t *testing.T) {
	repositoryTest(t, func(db *sql.DB, t *testing.T) {
		appRepo := NewAPPRepository(db)
		repo := NewClientRepository(db)

		// create app because it is foreign key that bounds these together
		if err := appRepo.Save(entity.Application{
			RecordID:    "test_app_id",
			Name:        "test_name",
			Description: "test_desc",
			Status:      "Active",
			CreateTime:  time.Now().UTC().Round(time.Second),
			CancelTime:  nil,
		}); err != nil {
			t.Errorf("application save operation failed on test prep process, err: %v", err)
		}

		client := entity.Client{
			RecordID:             "test_id",
			ClientID:             "test_client_id",
			APPID:                "test_app_id",
			Status:               "Active",
			RegisterTime:         time.Now().UTC().Round(time.Second),
			LastStatusChangeTime: time.Now().UTC().Round(time.Second),
		}

		if err := repo.Save(client); err != nil {
			t.Errorf("error occurred while saving, err: %v", err)
		}

		clientFromDB := entity.Client{}
		query := "SELECT RECORD_ID, CLIENT_ID, APP_ID, STATUS, REGISTER_TIME, LAST_STATUS_CHANGE_TIME FROM C_APP_USER WHERE RECORD_ID=?"
		if err := db.QueryRow(query, "test_id").Scan(
			&clientFromDB.RecordID,
			&clientFromDB.ClientID,
			&clientFromDB.APPID,
			&clientFromDB.Status,
			&clientFromDB.RegisterTime,
			&clientFromDB.LastStatusChangeTime,
		); err != nil {
			t.Errorf("test query failed, err: %v", err)
		}

		if !reflect.DeepEqual(clientFromDB, client) {
			t.Errorf("expected and given client are not deeply equal, expected: %v, \nfromDb: %v", client, clientFromDB)
		}
	})
}
