package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	hostname = "localhost"
	port     = 9999
)

type Environment struct {
	appService APPService
}

var env *Environment = nil

func initDB() {
	db, err := sql.Open("mysql", "bilginyuksel:toor@tcp(127.0.0.1:3306)/notificationservice")

	if err != nil {
		panic(err.Error())
	}

	appRepo := NewAPPRepository(db)

	env = &Environment{
		appService: NewAPPService(appRepo),
	}

	log.Printf("db connection established")
}

func main() {
	go initDB()
	StartWithHostAndPort(hostname, port)
}

func StartWithHostAndPort(hostname string, port int) {
	http.HandleFunc("/api/notification", pushNotificationClient)
	http.HandleFunc("/api/application", createNewApplication)

	url := fmt.Sprintf("%s:%d", hostname, port)

	log.Printf("server up and running, url: %v", url)
	log.Fatal(http.ListenAndServe(url, nil))

}

func pushNotificationClient(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	clientId := params.Get("clientId")

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("request body couldn't converted to bytes, err: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("internal server error"))
		return
	}

	notificationReq := &NotificationRequest{}
	err = json.Unmarshal(bytes, notificationReq)

	if err != nil {
		log.Println("illegal argument detected. notification req is incorrect. current body: ", string(bytes))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("illegal argument"))
		return
	}

	go pushNotificationToClient(clientId, *notificationReq)
}

func createNewApplication(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("request couldn't converted")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	createAppReq := &CreateAppRequest{}
	if err = json.Unmarshal(bytes, createAppReq); err != nil {
		log.Printf("illegal argument exception, request body is incorrect, err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("illegal argument"))
		return
	}

	if app, err := env.appService.CreateNewAPP(*createAppReq); err != nil {
		log.Printf("internal server error, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	} else {
		bytes, _ := json.Marshal(app)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	}

}
