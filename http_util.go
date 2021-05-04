package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bilginyuksel/push-notification/response"
)

// writeJSON send http response immediately
func writeJSON(w http.ResponseWriter, resp interface{}, httpStatus int) {
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(resp)

	if err != nil {
		resp := response.BaseResponse{
			ReturnCode: "9999",
			ReturnDesc: "Internal Server Error",
		}

		b, _ := json.Marshal(resp)

		log.Printf("interface could not converted to byte, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
	} else {
		w.WriteHeader(httpStatus)
		w.Write(bytes)
	}
}
