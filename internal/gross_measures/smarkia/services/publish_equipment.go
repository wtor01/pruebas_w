package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type PublishEquipmentDto struct {
	ProcessName     string
	DistributorId   string
	SmarkiaId       string
	CtId            string
	DistributorCDOS string
	Date            time.Time
}

type PublishEquipmentService struct {
	api       smarkia.GetEquipmentser
	publisher event.PublisherCreator
	topic     string
}

func NewPublishEquipmentService(publisher event.PublisherCreator, api smarkia.GetEquipmentser, topic string) *PublishEquipmentService {
	return &PublishEquipmentService{api: api, publisher: publisher, topic: topic}
}

func (svc PublishEquipmentService) Handle(ctx context.Context, dto PublishEquipmentDto) error {
	publisher, err := svc.publisher(ctx)
	if err != nil {
		return err
	}

	defer publisher.Close()

	equipments, err := svc.api.GetEquipments(ctx, smarkia.GetEquipmentsQuery{
		Id: dto.CtId,
	})
	if err != nil {
		return err
	}
	messages := make([]smarkia.EquipmentProcessEvent, 0, cap(equipments))

	for _, e := range equipments {
		messages = append(messages, smarkia.NewEquipmentProcessEvent(smarkia.MessageEquipmentDto{
			ProcessName:     dto.ProcessName,
			DistributorId:   dto.DistributorId,
			SmarkiaId:       dto.SmarkiaId,
			CtId:            dto.CtId,
			DistributorCDOS: dto.DistributorCDOS,
			Date:            dto.Date,
		}, e.ID, e.CUPS))
	}

	err = event.PublishAllEvents(ctx, svc.topic, svc.publisher, messages)

	return err
}
