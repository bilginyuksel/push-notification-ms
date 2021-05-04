package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterNotificationEndpoints(router *mux.Router) {
	router.HandleFunc("/api/notification", PushNotification).Methods("POST")
}

func PushNotification(w http.ResponseWriter, r *http.Request) {

}
