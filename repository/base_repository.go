package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// bilginyuksel:toor@tcp(127.0.0.1:3306)/notificationservice
	DefaultDBConn = DBConn{
		Provider: "mysql",
		Username: "bilginyuksel",
		Password: "toor",
		ConnType: "tcp",
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "notificationservice",
	}
)

type DBConn struct {
	Provider string
	Username string
	Password string
	Host     string
	ConnType string
	Port     int
	DBName   string
}

func (conn DBConn) buildMysqlDataSourceName() string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true",
		conn.Username,
		conn.Password,
		conn.ConnType,
		conn.Host,
		conn.Port,
		conn.DBName)
}

func ConnectMySQL(conn DBConn) (*sql.DB, error) {
	db, err := sql.Open(conn.Provider, conn.buildMysqlDataSourceName())

	if err != nil {
		fmt.Printf("db connection failed, err: %v", err)
		return nil, err
	}

	fmt.Println("db connection established")

	return db, nil
}
