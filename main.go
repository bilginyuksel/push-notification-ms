package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bilginyuksel/push-notification/entity"
	"github.com/bilginyuksel/push-notification/repository"
	"github.com/bilginyuksel/push-notification/request"
	"github.com/bilginyuksel/push-notification/service"
	"github.com/gorilla/mux"
)

const (
	hostname = "localhost"
	port     = 9999
)

type Environment struct {
	appService          service.APPService
	topicService        service.TopicService
	clientService       service.ClientService
	notificationService service.NotificationService
}

var env *Environment = nil

func initEnv() {

	db, err := repository.ConnectMySQL(repository.DefaultDBConn)

	if err != nil {
		panic(err)
	}

	appRepo := repository.NewAPPRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	clientRepo := repository.NewClientRepository(db)

	appService := service.NewAPPService(appRepo)
	topicService := service.NewTopicService(topicRepo)
	clientService := service.NewClientService(clientRepo, appService)
	notificationService := service.NewNotificationService()

	env = &Environment{
		appService:          appService,
		topicService:        topicService,
		clientService:       clientService,
		notificationService: notificationService,
	}
}

func main() {
	go initEnv()

	router := mux.NewRouter()

	appRouter := router.PathPrefix("/application").Subrouter()
	notificationRouter := router.PathPrefix("/notification").Subrouter()
	topicRouter := router.PathPrefix("/topic").Subrouter()
	clientRouter := router.PathPrefix("/client").Subrouter()

	RegisterApplicationEndpoints(appRouter)
	RegisterNotificationEndpoints(notificationRouter)
	RegisterTopicEndpoints(topicRouter)
	RegisterClientEndpoints(clientRouter)

	http.Handle("/", router)

	url := fmt.Sprintf("%s:%d", hostname, port)

	log.Printf("server up and running, url: %v", url)

	log.Fatal(http.ListenAndServe(url, nil))
	// StartWithHostAndPort(hostname, port)
}

func StartWithHostAndPort(hostname string, port int) {
	http.HandleFunc("/api/notification", pushNotificationClient)
	http.HandleFunc("/api/application", createNewApplication)
	http.HandleFunc("/api/topic", createNewApplication)
	http.HandleFunc("/api/client", createNewApplication)

	url := fmt.Sprintf("%s:%d", hostname, port)

	log.Printf("server up and running, url: %v", url)
	log.Fatal(http.ListenAndServe(url, nil))

}

func pushNotificationClient(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	clientId := params.Get("clientId")

	log.Println(clientId)

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("request body couldn't converted to bytes, err: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("internal server error"))
		return
	}

	notificationReq := &request.NotificationRequest{}
	err = json.Unmarshal(bytes, notificationReq)

	if err != nil {
		log.Println("illegal argument detected. notification req is incorrect. current body: ", string(bytes))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("illegal argument"))
		return
	}

	notif := entity.Notification{
		Title:   notificationReq.Title,
		Message: notificationReq.Message,
		Extras:  notificationReq.Extras,
	}
	go env.notificationService.Push(notif)
	// TODO: send notification via messaging queue
	// go pushNotificationToClient(clientId, *notificationReq)
}

func createNewApplication(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("request couldn't converted")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	createAppReq := &request.CreateAppRequest{}
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
