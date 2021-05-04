package main

import (
	"io/ioutil"
	"net/http"

	"github.com/bilginyuksel/push-notification/request"
	"github.com/bilginyuksel/push-notification/response"
	"github.com/bilginyuksel/push-notification/validation"
	"github.com/gorilla/mux"
)

var (
	OkResponse = response.BaseResponse{
		ReturnCode: "0",
		ReturnDesc: "OK",
	}

	IllegalArgumentException = response.BaseResponse{
		ReturnCode: "00012",
		ReturnDesc: "Illegal Argument Exception",
	}

	ServiceError = response.BaseResponse{
		ReturnCode: "00013",
		ReturnDesc: "Service error",
	}
)

func RegisterApplicationEndpoints(router *mux.Router) {
	router.HandleFunc("/create", CreateApplication).Methods("POST")
	router.HandleFunc("/query", QueryApplication).Methods("GET")
	router.HandleFunc("/query-by-id", QueryApplicationByID).Methods("GET")
	router.HandleFunc("/delete-by-id", DeleteApplication).Methods("DELETE")
}

func CreateApplication(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)

	req := &request.CreateAppRequest{}
	isValid := validation.Validate(b, req)

	if !isValid {
		writeJSON(w, IllegalArgumentException, 402)
		return
	}

	app, err := env.appService.CreateNewAPP(*req)

	if err != nil {
		writeJSON(w, ServiceError, 403)
		return
	}

	writeJSON(w, response.CreateAppResponse{
		AppID:        app.RecordID,
		BaseResponse: OkResponse,
	}, 200)
}

func QueryApplication(w http.ResponseWriter, r *http.Request) {
	apps, err := env.appService.GetAll()

	if err != nil {
		writeJSON(w, ServiceError, 403)
		return
	}

	writeJSON(w, response.QueryAppResponse{
		BaseResponse: OkResponse,
		Apps:         apps,
	}, 200)
}

func QueryApplicationByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func UpdateApplication(w http.ResponseWriter, r *http.Request) {

}

func DeleteApplication(w http.ResponseWriter, r *http.Request) {

}
