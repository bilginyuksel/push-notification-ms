package main

type PushNotificationReq struct {
	Title   string            `json:"title"`
	Message string            `json:"message"`
	Extras  map[string]string `json:"extras"`
}
