package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/bilginyuksel/push-notification/entity"
	"github.com/hashicorp/go-uuid"
)

func TestAppRepositorySave(t *testing.T) {

	db := prepareTestDB()
	repo := NewAPPRepository(db)

	uniqueRecordId, _ := uuid.GenerateUUID()

	application := entity.Application{
		RecordID:    uniqueRecordId,
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
}
