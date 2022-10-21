package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ListFeaturesDto struct {
	Limit  int
	Offset *int
}
type ListFeatures struct {
	repository aggregations.AggregationsFeaturesRepository
	tracer     trace.Tracer
}

func NewListFeatures(repository aggregations.AggregationsFeaturesRepository) *ListFeatures {
	return &ListFeatures{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

func (f ListFeatures) Handler(ctx context.Context, dto ListFeaturesDto) ([]aggregations.Features, int, error) {
	ctx, span := f.tracer.Start(ctx, "SearchFiscalBillingMeasures - Handler")

	result, count, err := f.repository.ListFeatures(ctx, db.Pagination{
		Limit:  dto.Limit,
		Offset: dto.Offset,
	})
	if err != nil {
		return []aggregations.Features{}, 0, err
	}

	span.SetAttributes(attribute.Int("response", len(result)))
	return result, count, err
}
