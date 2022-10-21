package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type GetAggregations struct {
	Location   *time.Location
	repository aggregations.AggregationMongoRepository
	tracer     trace.Tracer
}

func NewGetAggregationsDashboard(repository aggregations.AggregationMongoRepository, location *time.Location) *GetAggregations {
	return &GetAggregations{
		repository: repository,
		Location:   location,
		tracer:     telemetry.GetTracer(),
	}
}

func (s GetAggregations) Handler(ctx context.Context, aggregationDto aggregations.GetAggregationsDto) ([]aggregations.Aggregation, int64, error) {
	ctx, span := s.tracer.Start(ctx, "GetAggregations - Handler")
	defer span.End()

	ags, count, err := s.repository.GetAggregations(ctx, aggregationDto)

	if err != nil {
		return []aggregations.Aggregation{}, 0, err
	}

	return ags, count, nil
}
