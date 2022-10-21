package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
)

type Services struct {
	PublishDistributorService         *PublishDistributorService
	PublishCtService                  *PublishCtService
	PublishEquipmentService           *PublishEquipmentService
	PublishMeasuresFromSmarkiaService *PublishMeasuresFromSmarkiaService
}

type ApiSmarkia interface {
	smarkia.GetCTer
	smarkia.GetEquipmentser
	smarkia.GetMagnitudeser
	smarkia.GetClosinger
}

func NewServices(
	publisher event.PublisherCreator,
	clientInventory clients.Inventory,
	api ApiSmarkia,
	cnf config.Config,
	inventoryRepository measures.InventoryRepository,
) *Services {

	publishDistributorService := NewPublishDistributorService(publisher, clientInventory, cnf.SmarkiaTopic)
	publishCtService := NewPublishCtService(publisher, api, cnf.SmarkiaTopic)
	publishEquipmentService := NewPublishEquipmentService(publisher, api, cnf.SmarkiaTopic)
	publishMeasuresFromSmarkiaService := NewPublishMeasuresFromSmarkiaService(api, publisher, cnf.TopicMeasures, inventoryRepository)

	return &Services{
		PublishDistributorService:         publishDistributorService,
		PublishCtService:                  publishCtService,
		PublishEquipmentService:           publishEquipmentService,
		PublishMeasuresFromSmarkiaService: publishMeasuresFromSmarkiaService,
	}
}
