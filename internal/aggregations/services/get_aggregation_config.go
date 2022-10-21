package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/trace"
)

type GetAggregationConfigsServiceDto struct {
	Q      string
	Limit  int
	Offset *int
}

func NewGetAggregationConfigsServiceDto(q string, limit int, offset *int) GetAggregationConfigsServiceDto {
	return GetAggregationConfigsServiceDto{
		Q:      q,
		Limit:  limit,
		Offset: offset,
	}
}

type GetAggregationConfigsService struct {
	aggregationRepository aggregations.AggregationConfigRepository
	tracer                trace.Tracer
}

func NewGetAggregationConfigService(aggregationRepository aggregations.AggregationConfigRepository) *GetAggregationConfigsService {
	return &GetAggregationConfigsService{
		aggregationRepository: aggregationRepository,
		tracer:                telemetry.GetTracer(),
	}
}

func (s GetAggregationConfigsService) Handler(ctx context.Context, dto GetAggregationConfigsServiceDto) ([]aggregations.Config, int, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllAggregations - Handler")
	defer span.End()
	configs, count, err := s.aggregationRepository.GetAggregationConfigs(ctx, aggregations.GetConfigsQuery{
		Query:  dto.Q,
		Limit:  dto.Limit,
		Offset: dto.Offset,
	})

	return configs, count, err
}
