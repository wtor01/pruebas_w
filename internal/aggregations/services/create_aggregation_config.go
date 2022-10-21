package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type CreateAggregationConfigService struct {
	aggregationRepository aggregations.AggregationConfigRepository
	schedulerCreator      scheduler.ClientCreator
	topic                 string
	tracer                trace.Tracer
	location              *time.Location
	featuresRepository    aggregations.AggregationsFeaturesRepository
}

type CreateConfigDto struct {
	Name        string
	Description string
	Scheduler   string
	StartDate   time.Time
	EndDate     time.Time
	Features    []aggregations.ConfigFeatureDto
}

func NewCreateConfigDto(name, scheduler string, description *string, startDate time.Time, endDate *time.Time, features []aggregations.ConfigFeatureDto) CreateConfigDto {

	configDto := CreateConfigDto{
		Name:      name,
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

func NewCreateAggregationConfigService(repository aggregations.AggregationConfigRepository, featuresRepository aggregations.AggregationsFeaturesRepository, schedulerCreator scheduler.ClientCreator, topic string, loc *time.Location) *CreateAggregationConfigService {
	return &CreateAggregationConfigService{
		aggregationRepository: repository,
		schedulerCreator:      schedulerCreator,
		topic:                 topic,
		tracer:                telemetry.GetTracer(),
		location:              loc,
		featuresRepository:    featuresRepository,
	}
}

func (s CreateAggregationConfigService) Handler(ctx context.Context, aggregationConfigDto CreateConfigDto) (aggregations.Config, error) {
	ctx, span := s.tracer.Start(ctx, "CreateAggregationConfig - Handler")
	defer span.End()

	features, err := s.featuresRepository.GetFeaturesByIds(ctx, utils.MapSlice(aggregationConfigDto.Features, func(item aggregations.ConfigFeatureDto) string {
		return item.Id
	}))
	if err != nil {
		return aggregations.Config{}, err
	}

	aggregationConfig, err := aggregations.NewConfig(
		aggregationConfigDto.Name,
		aggregationConfigDto.Description,
		aggregationConfigDto.Scheduler,
		aggregationConfigDto.StartDate.In(s.location).UTC(),
		aggregationConfigDto.EndDate.In(s.location).UTC(),
		features,
	)

	if err != nil {
		return aggregations.Config{}, err
	}

	schedulerClient, err := s.schedulerCreator(ctx)
	if err != nil {
		return aggregations.Config{}, err
	}
	defer schedulerClient.Close()

	jobId, err := schedulerClient.CreateJob(ctx, aggregationConfig, s.topic)
	if err != nil {
		return aggregations.Config{}, err
	}
	aggregationConfig.SetSchedulerId(jobId)

	result, err := s.aggregationRepository.SaveAggregationConfig(ctx, aggregationConfig)
	if err != nil {
		schedulerClient.DeleteJob(ctx, aggregationConfig.SchedulerId)
		return aggregations.Config{}, err
	}
	return result, nil
}
