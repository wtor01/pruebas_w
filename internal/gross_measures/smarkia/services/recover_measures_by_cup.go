package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

type RecoverSmarkiaMeasuresDTO struct {
	CUPS          string
	DistributorID string
	Date          time.Time
	ProcessName   smarkia.ProcessName
}

type Api interface {
	smarkia.GetMagnitudeser
	smarkia.GetClosinger
	smarkia.GetEquipmentser
}

type RecoverSmarkiaMeasures struct {
	api                               smarkia.GetEquipmentser
	inventoryClient                   clients.Inventory
	publishMeasuresFromSmarkiaService *PublishMeasuresFromSmarkiaService
	Location                          *time.Location
}

func NewRecoverSmarkiaMeasures(
	api Api,
	inventoryClient clients.Inventory,
	location *time.Location,
	publisher event.PublisherCreator,
	topicMeasure string,
	inventoryRepository measures.InventoryRepository,
) *RecoverSmarkiaMeasures {
	return &RecoverSmarkiaMeasures{
		api:             api,
		inventoryClient: inventoryClient,
		Location:        location,
		publishMeasuresFromSmarkiaService: NewPublishMeasuresFromSmarkiaService(
			api,
			publisher,
			topicMeasure,
			inventoryRepository,
		),
	}
}

func (svc RecoverSmarkiaMeasures) Handler(ctx context.Context, dto RecoverSmarkiaMeasuresDTO) error {

	distributor, err := svc.inventoryClient.GetDistributorById(ctx, dto.DistributorID)

	if err != nil {
		return err
	}

	equipments, err := svc.api.GetEquipments(ctx, smarkia.GetEquipmentsQuery{
		Cups: dto.CUPS,
	})

	if err != nil {
		return err
	}

	if len(equipments) == 0 {
		return errors.New(fmt.Sprintf("invalid cup %s", dto.CUPS))
	}

	equipment := equipments[0]

	ev := smarkia.NewEquipmentProcessEvent(smarkia.MessageEquipmentDto{
		ProcessName:     dto.ProcessName,
		DistributorId:   distributor.ID,
		SmarkiaId:       distributor.SmarkiaID,
		CtId:            equipment.CtId,
		DistributorCDOS: distributor.CDOS,
		Date:            utils.ToUTC(dto.Date.AddDate(0, 0, 1), svc.Location),
	},
		equipment.ID,
		equipment.CUPS,
	)

	err = svc.publishMeasuresFromSmarkiaService.Handle(ctx, ev)

	return err
}
