package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterTopicEndpoints(router *mux.Router) {
	router.HandleFunc("/api/topic", CreateTopic).Methods("POST")
	router.HandleFunc("/api/topic", QueryTopic).Methods("GET")
	router.HandleFunc("/api/topic", UpdateTopic).Methods("PUT")
	router.HandleFunc("/api/topic", DeleteTopic).Methods("DELETE")
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {

}

func QueryTopic(w http.ResponseWriter, r *http.Request) {

}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {

}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {

}
