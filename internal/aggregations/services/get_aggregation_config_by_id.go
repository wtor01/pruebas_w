package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetAggregationConfigByIdService struct {
	aggregationRepository aggregations.AggregationConfigRepository
	tracer                trace.Tracer
}

func NewGetAggregationConfigByIdService(repository aggregations.AggregationConfigRepository) *GetAggregationConfigByIdService {
	return &GetAggregationConfigByIdService{
		aggregationRepository: repository,
		tracer:                telemetry.GetTracer(),
	}
}

func (s GetAggregationConfigByIdService) Handler(ctx context.Context, aggregationConfigId string) (aggregations.Config, error) {
	ctx, span := s.tracer.Start(ctx, "GetAggregationConfigById - Handler")
	defer span.End()

	span.SetAttributes(attribute.String("aggregation_config_id", aggregationConfigId))
	result, err := s.aggregationRepository.GetAggregationConfigById(ctx, aggregationConfigId)

	if err != nil {
		return aggregations.Config{}, err
	}

	return result, err
}
