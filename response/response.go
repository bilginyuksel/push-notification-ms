package response

import "github.com/bilginyuksel/push-notification/entity"

type BaseResponse struct {
	ReturnCode string `json:"returnCode"`
	ReturnDesc string `json:"returnDesc"`
}

type CreateAppResponse struct {
	BaseResponse
	AppID string `json:"appId"`
}

type QueryAppResponse struct {
	BaseResponse
	Apps []*entity.Application `json:"apps"`
}
