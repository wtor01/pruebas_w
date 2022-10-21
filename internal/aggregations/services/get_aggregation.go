package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type GetAggregation struct {
	Location   *time.Location
	repository aggregations.AggregationMongoRepository
	tracer     trace.Tracer
}

func NewGetAggregation(repository aggregations.AggregationMongoRepository, location *time.Location) *GetAggregation {
	return &GetAggregation{
		repository: repository,
		Location:   location,
		tracer:     telemetry.GetTracer(),
	}
}

func (s GetAggregation) Handler(ctx context.Context, aggregationDto aggregations.GetAggregationDto) (aggregations.AggregationPrevious, error) {
	ctx, span := s.tracer.Start(ctx, "GetAggregation - Handler")
	defer span.End()

	ag, err := s.repository.GetPreviousAggregation(ctx, aggregationDto)

	if err != nil {
		return aggregations.AggregationPrevious{}, err
	}

	return ag, nil
}
