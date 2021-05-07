package request

type NotificationRequest struct {
	Title   string            `json:"title"`
	Message string            `json:"message"`
	Extras  map[string]string `json:"extras"`
}

type CreateAppRequest struct {
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateTopicRequest struct {
	Application string `json:"application"`
	Topic       string `json:"topic"`
}

type GetClientIDRequest struct {
	ApplicationID string `json:"applicationId"`
}

type CreateClientRequest struct {
	UserID   string `json:"userId"`
	APPID    string `json:"appId"`
	ClientID string `json:"clientId"`
}

type PushNotificationToUserRequest struct {
	Application  string               `json:"application"`
	ClientID     string               `json:"clientId"`
	Notification *NotificationRequest `json:"notification"`
}

type PushNotificationByTopicRequest struct {
	Application  string               `json:"application"`
	Topic        string               `json:"topic"`
	Notification *NotificationRequest `json:"notification"`
}
