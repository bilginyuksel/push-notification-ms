package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	hostname = "localhost"
	port     = 9999
)

func main() {
	StartWithHostAndPort(hostname, port)
}

func StartWithHostAndPort(hostname string, port int) {
	http.HandleFunc("/api/notification", pushNotificationClient)

	url := fmt.Sprintf("%s:%d", hostname, port)

	log.Fatal(http.ListenAndServe(url, nil))
}

func pushNotificationClient(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	clientId := params.Get("clientId")

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("request body couldn't converted to bytes, err: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Internal server error"))
		return
	}

	notificationReq := &PushNotificationReq{}
	err = json.Unmarshal(bytes, notificationReq)

	if err != nil {
		log.Println("Illegal argument detected. Notification Req is incorrect. Current body: ", string(bytes))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Illegal Argument"))
		return
	}

	go func(clientId string, notification PushNotificationReq) {
		if err := pushNotificationToClient(clientId, notification); err != nil {
			log.Println("Individual message sent error, err: ", err)
		}
	}(clientId, *notificationReq)
}
