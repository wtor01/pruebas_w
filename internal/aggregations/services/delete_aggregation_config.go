package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DeleteAggregationConfigService struct {
	aggregationRepository aggregations.AggregationConfigRepository
	schedulerCreator      scheduler.ClientCreator
	topic                 string
	tracer                trace.Tracer
}

func NewDeleteAggregationConfigService(repository aggregations.AggregationConfigRepository, schedulerCreator scheduler.ClientCreator, topic string) *DeleteAggregationConfigService {
	return &DeleteAggregationConfigService{
		aggregationRepository: repository,
		schedulerCreator:      schedulerCreator,
		topic:                 topic,
		tracer:                telemetry.GetTracer(),
	}
}

func (s DeleteAggregationConfigService) Handler(ctx context.Context, aggregationConfigId string) error {
	ctx, span := s.tracer.Start(ctx, "DeleteAggregationConfig - Handler")
	defer span.End()

	span.SetAttributes(attribute.String("AggregationConfigId", aggregationConfigId))
	aggregationConfig, err := s.aggregationRepository.GetAggregationConfigById(ctx, aggregationConfigId)

	if err != nil {
		return err
	}

	schedulerClient, err := s.schedulerCreator(ctx)
	if err != nil {
		return err
	}
	defer schedulerClient.Close()
	
	err = schedulerClient.DeleteJob(ctx, aggregationConfig.SchedulerId)
	if err != nil {
		return err
	}

	err = s.aggregationRepository.DeleteAggregationConfig(ctx, aggregationConfig.Id)

	if err != nil {
		schedulerClient.CreateJob(ctx, aggregationConfig, s.topic)
		return err
	}
	return err
}
