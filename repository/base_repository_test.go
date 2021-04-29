package repository

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"
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
	// insertDummyData(db)
	return db
}

func destroyDB(db *sql.DB) {
	db.Exec("DROP DATABASE testnotificationservice;")
	log.Printf("test db destroyed")
}
