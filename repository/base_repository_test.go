package repository

import (
	"database/sql"
	"io/ioutil"
	"testing"
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

	_, err = db.Exec(string(file))

	if err != nil {
		panic(err)
	}
}

func insertDummyData(db *sql.DB) {

}

func prepareTestDB() *sql.DB {
	db := createTestDB()
	createTestDBSchema(db)
	insertDummyData(db)
	return db
}
