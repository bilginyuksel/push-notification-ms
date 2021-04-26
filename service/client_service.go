package service

import (
	"errors"
	"log"
	"time"

	"github.com/bilginyuksel/push-notification/entity"
	"github.com/bilginyuksel/push-notification/repository"
	"github.com/bilginyuksel/push-notification/request"
	"github.com/hashicorp/go-uuid"
)

type clientServiceImpl struct {
	clientRepo repository.ClientRepository
	appService APPService
}

type ClientService interface {
	CreateNewClient(req request.CreateClientRequest) (*entity.Client, error)
}

func NewClientService(repo repository.ClientRepository, appService APPService) ClientService {
	return &clientServiceImpl{
		clientRepo: repo,
		appService: appService}
}

func (self *clientServiceImpl) CreateNewClient(req request.CreateClientRequest) (*entity.Client, error) {

	if isExist := self.appService.IsExist(req.APPID); !isExist {
		log.Println("app could not found")
		return nil, errors.New("app could not found")
	}

	uuid, _ := uuid.GenerateUUID()

	client := entity.Client{
		RecordID:             uuid,
		APPID:                req.APPID,
		ClientID:             req.ClientID,
		RegisterTime:         time.Now(),
		LastStatusChangeTime: time.Now(),
		Status:               "Approving",
	}

	if err := self.clientRepo.Save(client); err != nil {
		log.Printf("an error occured while saving client, err: %v", err)
		return nil, err
	}

	log.Printf("client saved successfully, client: %v", client)

	return &client, nil
}
