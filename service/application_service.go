package service

import (
	"log"
	"time"

	"github.com/hashicorp/go-uuid"

	"github.com/bilginyuksel/push-notification/entity"
	"github.com/bilginyuksel/push-notification/repository"
	"github.com/bilginyuksel/push-notification/request"
)

type appServiceImpl struct {
	repo repository.APPRepository
}

type APPService interface {
	CreateNewAPP(req request.CreateAppRequest) (*entity.Application, error)
	GetAll() ([]*entity.Application, error)
	IsExist(appID string) bool
}

func NewAPPService(appRepo repository.APPRepository) APPService {
	return &appServiceImpl{
		repo: appRepo,
	}
}

func (service *appServiceImpl) CreateNewAPP(req request.CreateAppRequest) (*entity.Application, error) {
	recordID := ""

	if uuid, err := uuid.GenerateUUID(); err != nil {
		log.Printf("uuid generation failed, err: %v", err)
	} else {
		recordID = uuid
	}

	app := entity.Application{
		RecordID:    recordID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "Approving",
		CreateTime:  time.Now(),
		CancelTime:  nil,
	}

	if err := service.repo.Save(app); err != nil {
		log.Printf("application couldn't created, err: %v", err)
		return nil, err
	}

	log.Printf("new application created, app: %v", app)
	return &app, nil
}

func (service *appServiceImpl) IsExist(appID string) bool {
	return service.repo.IsExist(appID)
}

func (service *appServiceImpl) GetAll() ([]*entity.Application, error) {
	return service.repo.GetAll()
}
