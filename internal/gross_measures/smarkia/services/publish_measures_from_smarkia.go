package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type ApiMeasures interface {
	smarkia.GetMagnitudeser
	smarkia.GetClosinger
}

type PublishMeasuresFromSmarkiaService struct {
	api                 ApiMeasures
	publisher           event.PublisherCreator
	topicMeasure        string
	inventoryRepository measures.InventoryRepository
}

func NewPublishMeasuresFromSmarkiaService(
	api ApiMeasures,
	publisher event.PublisherCreator,
	topicMeasure string,
	inventoryRepository measures.InventoryRepository,
) *PublishMeasuresFromSmarkiaService {
	return &PublishMeasuresFromSmarkiaService{
		api:                 api,
		publisher:           publisher,
		topicMeasure:        topicMeasure,
		inventoryRepository: inventoryRepository,
	}
}

func (svc PublishMeasuresFromSmarkiaService) handleCurve(ctx context.Context, dto smarkia.EquipmentProcessEvent) error {
	ms, err := svc.api.GetMagnitudes(ctx, dto.Payload.EquipmentId, dto.Payload.Date.UTC())

	if err != nil {
		return err
	}

	meter, _ := svc.inventoryRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        dto.Payload.CUPS,
		Time:        dto.Payload.Date.UTC(),
		Distributor: dto.Payload.DistributorId,
	})

	for i := range ms {
		ms[i].DistributorCDOS = dto.Payload.DistributorCDOS
		ms[i].DistributorID = dto.Payload.DistributorId
		ms[i].MeterSerialNumber = meter.SerialNumber()
	}

	events := gross_measures.ListMeasureCurveWriteToEvents(ms, gross_measures.MaxMeasuresInEvent)

	err = event.PublishAllEvents(ctx, svc.topicMeasure, svc.publisher, events)

	return err
}

func (svc PublishMeasuresFromSmarkiaService) handleClose(ctx context.Context, dto smarkia.EquipmentProcessEvent) error {
	ms, err := svc.api.GetClosinger(ctx, dto.Payload.EquipmentId, dto.Payload.Date.UTC())

	if err != nil {
		return err
	}

	meter, _ := svc.inventoryRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        dto.Payload.CUPS,
		Time:        dto.Payload.Date.UTC(),
		Distributor: dto.Payload.DistributorId,
	})

	for i := range ms {
		ms[i].DistributorCDOS = dto.Payload.DistributorCDOS
		ms[i].DistributorID = dto.Payload.DistributorId
		ms[i].MeterSerialNumber = meter.SerialNumber()
	}

	events := gross_measures.ListMeasureCloseWriteToEvents(ms, gross_measures.MaxMeasuresInEvent)

	err = event.PublishAllEvents(ctx, svc.topicMeasure, svc.publisher, events)

	return err
}

func (svc PublishMeasuresFromSmarkiaService) Handle(ctx context.Context, dto smarkia.EquipmentProcessEvent) error {
	if dto.Payload.ProcessName == smarkia.ProcessCurve {
		return svc.handleCurve(ctx, dto)
	}
	return svc.handleClose(ctx, dto)
}
