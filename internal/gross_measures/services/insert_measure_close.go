package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type InsertMeasureCloseService struct {
	InsertMeasureBase
	measureRepository gross_measures.GrossMeasureRepository
	publisher         event.PublisherCreator
}

func NewInsertMeasureCloseService(
	measureRepository gross_measures.GrossMeasureRepository,
	validationClient clients.Validation,
	inventoryClient clients.Inventory,
	publisher event.PublisherCreator,
	topic string,
) *InsertMeasureCloseService {
	return &InsertMeasureCloseService{
		InsertMeasureBase: InsertMeasureBase{
			validationClient: validationClient,
			inventoryClient:  inventoryClient,
			topic:            topic,
			generatorDate:    time.Now,
		},
		publisher:         publisher,
		measureRepository: measureRepository,
	}
}

func (svc InsertMeasureCloseService) Handle(ctx context.Context, measuresClose []gross_measures.MeasureCloseWrite) error {
	measuresBase := make([]gross_measures.GrossMeasureBase, 0, cap(measuresClose))

	for i := range measuresClose {
		measuresClose[i].GenerationDate = svc.generatorDate().UTC()
		measuresBase = append(measuresBase, &measuresClose[i])
	}

	err := svc.SetInsertMeasureMetadata(ctx, measuresBase)

	if err != nil {
		return err
	}

	err = svc.ValidateMeasure(ctx, measuresBase)

	if err != nil {
		return err
	}

	err = svc.measureRepository.SaveAllMeasuresClose(ctx, measuresClose)

	return err
}
