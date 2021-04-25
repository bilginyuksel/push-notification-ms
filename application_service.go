package main

import (
	"log"
	"time"

	"github.com/hashicorp/go-uuid"
)

type appServiceImpl struct {
	repo APPRepository
}

type APPService interface {
	CreateNewAPP(req CreateAppRequest) (*Application, error)
	IsExist(appID string) bool
}

func NewAPPService(appRepo APPRepository) APPService {
	return &appServiceImpl{
		repo: appRepo,
	}
}

func (service *appServiceImpl) CreateNewAPP(req CreateAppRequest) (*Application, error) {
	// generate recordID

	recordID := ""
	if uuid, err := uuid.GenerateUUID(); err != nil {
		log.Printf("uuid generation failed, err: %v", err)
	} else {
		recordID = uuid
	}

	app := Application{
		RecordID:    recordID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "Approving",
		CreateTime:  time.Now(),
		CancelTime:  nil,
	}

	if err := service.repo.save(app); err != nil {
		log.Printf("application couldn't created, err: %v", err)
		return nil, err
	}

	log.Printf("new application created, app: %v", app)
	return &app, nil
}

func (service *appServiceImpl) IsExist(appId string) bool {
	return false
}
