package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterClientEndpoints(router *mux.Router) {
	router.HandleFunc("/api/client", CreateClient).Methods("POST")
}

func CreateClient(w http.ResponseWriter, r *http.Request) {

}
