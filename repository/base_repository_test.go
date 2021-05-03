package repository

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/bilginyuksel/push-notification/entity"
)

func TestConnectMySQL(t *testing.T) {
	dbconn := DefaultDBConn
	dbconn.DBName = ""
	t.Logf("default db connection copied and db name changed, dbconn: %v", dbconn)

	if _, err := ConnectMySQL(dbconn); err != nil {
		t.Errorf("db connection failed, err: %v", err)
	}

	t.Log("db connection successfull")
}

func repositoryTest(t *testing.T, test func(db *sql.DB, t *testing.T)) {
	db := prepareTestDB()
	test(db, t)
	destroyDB(db)
}

func createTestDB() *sql.DB {
	dbconn := DefaultDBConn
	dbconn.DBName = ""

	db, _ := ConnectMySQL(dbconn)

	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS testnotificationservice")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE testnotificationservice")

	if err != nil {
		panic(err)
	}

	return db
}

func createTestDBSchema(db *sql.DB) {
	file, err := ioutil.ReadFile("../init.sql")

	if err != nil {
		panic(err)
	}

	queries := strings.Split(string(file), ";")
	length := len(queries)

	// don't take the last item because when it splits the sql query the last item is empty string
	for _, query := range queries[:length-1] {
		_, err = db.Exec(query)

		if err != nil {
			panic(err)
		}
	}
}

func prepareTestDB() *sql.DB {
	db := createTestDB()
	createTestDBSchema(db)
	return db
}

func destroyDB(db *sql.DB) {
	db.Exec("DROP DATABASE testnotificationservice;")
	log.Printf("test db destroyed")
}

type testPreperation struct {
	db *sql.DB
}

func (tp *testPreperation) createSampleApplication(recordId string) {
	appRepo := NewAPPRepository(tp.db)
	if err := appRepo.Save(entity.Application{
		RecordID:    recordId,
		Name:        "test_app_name",
		Description: "test_app_desc",
		Status:      "Active",
		CreateTime:  time.Now().UTC().Round(time.Second),
		CancelTime:  nil,
	}); err != nil {
		log.Fatalf("application save operation failed on test prep process, err: %v", err)
	}
}

func (tp *testPreperation) createSampleTopic(recordID, appID, name string) {
	topicRepo := NewTopicRepository(tp.db)

	if err := topicRepo.Save(entity.Topic{
		RecordID:    recordID,
		AppID:       appID,
		Name:        name,
		Description: "test_topic_desc",
	}); err != nil {
		log.Fatalf("topic save operation failed on test prep process, err: %v", err)
	}
}

func (tp *testPreperation) createSampleUser(recordId, appId string) {
	userRepo := NewClientRepository(tp.db)

	if err := userRepo.Save(entity.Client{
		RecordID:             recordId,
		ClientID:             "test_client_id",
		APPID:                appId,
		Status:               "Active",
		RegisterTime:         time.Now().UTC().Round(time.Second),
		LastStatusChangeTime: time.Now().UTC().Round(time.Second),
	}); err != nil {
		log.Fatalf("user save operation failed on test prep process, err: %v", err)
	}
}
