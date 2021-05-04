package response

type BaseResponse struct {
	ReturnCode string `json:"returnCode"`
	ReturnDesc string `json:"returnDesc"`
}

type CreateAppResponse struct {
	BaseResponse
	AppID string `json:"appId"`
}
