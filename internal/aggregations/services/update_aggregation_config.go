package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type UpdateAggregationConfigService struct {
	aggregationRepository aggregations.AggregationConfigRepository
	schedulerCreator      scheduler.ClientCreator
	topic                 string
	tracer                trace.Tracer
	location              *time.Location
	featuresRepository    aggregations.AggregationsFeaturesRepository
}

type UpdateConfigDto struct {
	Description string
	Scheduler   string
	StartDate   time.Time
	EndDate     time.Time
	Features    []aggregations.ConfigFeatureDto
}

func NewUpdateConfigDto(scheduler string, description *string, startDate time.Time, endDate *time.Time, features []aggregations.ConfigFeatureDto) UpdateConfigDto {

	configDto := UpdateConfigDto{
		Scheduler: scheduler,
		StartDate: startDate,
		Features:  features,
	}

	if endDate != nil {
		configDto.EndDate = *endDate
	}

	if description != nil {
		configDto.Description = *description
	}

	return configDto
}

func NewUpdateAggregationConfigService(repository aggregations.AggregationConfigRepository, featuresRepository aggregations.AggregationsFeaturesRepository, schedulerCreator scheduler.ClientCreator, topic string, loc *time.Location) *UpdateAggregationConfigService {
	return &UpdateAggregationConfigService{
		aggregationRepository: repository,
		schedulerCreator:      schedulerCreator,
		topic:                 topic,
		tracer:                telemetry.GetTracer(),
		location:              loc,
		featuresRepository:    featuresRepository,
	}

}

func (s UpdateAggregationConfigService) Handler(ctx context.Context, aggregationConfigId string, aggregationConfigDto UpdateConfigDto) (aggregations.Config, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateAggregationConfig - Handler")
	defer span.End()

	span.SetAttributes(attribute.String("AggregationConfigId", aggregationConfigId))

	schedulerClient, err := s.schedulerCreator(ctx)
	if err != nil {
		return aggregations.Config{}, err
	}
	defer schedulerClient.Close()

	aggregationConfig, err := s.aggregationRepository.GetAggregationConfigById(ctx, aggregationConfigId)

	if err != nil {
		return aggregations.Config{}, err
	}

	features, err := s.featuresRepository.GetFeaturesByIds(ctx, utils.MapSlice(aggregationConfigDto.Features, func(item aggregations.ConfigFeatureDto) string {
		return item.Id
	}))

	if err != nil {
		return aggregations.Config{}, err
	}

	oldAggregationConfig := aggregationConfig.Clone()

	err = aggregationConfig.Update(
		aggregationConfigDto.Scheduler,
		aggregationConfigDto.Description,
		aggregationConfigDto.StartDate.In(s.location).UTC(),
		aggregationConfigDto.EndDate.In(s.location).UTC(),
		features,
	)

	if err != nil {
		return aggregations.Config{}, err
	}

	err = schedulerClient.UpdateJob(ctx, aggregationConfig, s.topic)
	if err != nil {
		return aggregations.Config{}, err
	}

	result, err := s.aggregationRepository.SaveAggregationConfig(ctx, aggregationConfig)

	if err != nil {
		_ = schedulerClient.UpdateJob(ctx, oldAggregationConfig, s.topic)
		return aggregations.Config{}, err
	}

	return result, nil
}
