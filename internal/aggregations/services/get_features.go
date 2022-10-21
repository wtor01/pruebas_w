package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetFeatures struct {
	repository aggregations.AggregationsFeaturesRepository
	tracer     trace.Tracer
}

func NewGetFeaturesService(repository aggregations.AggregationsFeaturesRepository) *GetFeatures {
	return &GetFeatures{
		repository: repository,
		tracer:     telemetry.GetTracer()}
}

type GetFeaturesDTO struct {
	ID string
}

func (f GetFeatures) Handler(ctx context.Context, dto GetFeaturesDTO) (aggregations.Features, error) {
	if dto.ID == "" {
		return aggregations.Features{}, errors.New("Empty ID")
	}
	ctx, span := f.tracer.Start(ctx, "GetFeatures - Handler")
	defer span.End()

	feature, err := f.repository.GetFeatures(ctx, dto.ID)
	if err != nil {
		return aggregations.Features{}, err
	}
	span.SetAttributes(attribute.Bool("response", true))

	return feature, nil
}
