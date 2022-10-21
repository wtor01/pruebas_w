package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type InsertMeasureCurveService struct {
	InsertMeasureBase
	measureRepository gross_measures.GrossMeasureRepository
	publisher         event.PublisherCreator
}

func NewInsertMeasureCurveService(
	measureRepository gross_measures.GrossMeasureRepository,
	validationClient clients.Validation,
	inventoryClient clients.Inventory,
	publisher event.PublisherCreator,
	topic string,
) *InsertMeasureCurveService {
	return &InsertMeasureCurveService{
		InsertMeasureBase: InsertMeasureBase{
			validationClient: validationClient,
			inventoryClient:  inventoryClient,
			topic:            topic,
			generatorDate:    time.Now,
		},
		measureRepository: measureRepository,
		publisher:         publisher,
	}
}

func (svc InsertMeasureCurveService) Handle(ctx context.Context, measuresCurve []gross_measures.MeasureCurveWrite) error {

	measuresBase := make([]gross_measures.GrossMeasureBase, 0, cap(measuresCurve))

	for i := range measuresCurve {
		measuresCurve[i].GenerationDate = svc.generatorDate().UTC()
		measuresBase = append(measuresBase, &measuresCurve[i])
	}

	err := svc.SetInsertMeasureMetadata(ctx, measuresBase)

	if err != nil {
		return err
	}

	err = svc.ValidateMeasure(ctx, measuresBase)

	if err != nil {
		return err
	}

	err = svc.measureRepository.SaveAllMeasuresCurve(ctx, measuresCurve)

	return err
}
