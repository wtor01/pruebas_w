package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DeleteFeatures struct {
	repository aggregations.AggregationsFeaturesRepository
	tracer     trace.Tracer
}

func NewDeleteFeaturesService(repository aggregations.AggregationsFeaturesRepository) *DeleteFeatures {
	return &DeleteFeatures{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

type DeleteFeaturesDTO struct {
	ID string
}

func (f DeleteFeatures) Handler(ctx context.Context, dto DeleteFeaturesDTO) error {
	if dto.ID == "" {
		return errors.New("empty ID")
	}
	ctx, span := f.tracer.Start(ctx, "DeleteFeatures - Handler")
	defer span.End()

	features, err := f.repository.GetFeatures(ctx, dto.ID)

	if err != nil {
		return err
	}

	err = f.repository.DeleteFeatures(ctx, features.ID)
	span.SetAttributes(attribute.Bool("response", true))

	return err

}
