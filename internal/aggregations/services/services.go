package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"time"
)

type FeaturesServices struct {
	CreateFeatures *CreateFeatures
	GetFeatures    *GetFeatures
	DeleteFeatures *DeleteFeatures
	UpdateFeatures *UpdateFeatures
	ListFeatures   *ListFeatures
}

func NewFeaturesServices(repository aggregations.AggregationsFeaturesRepository) *FeaturesServices {
	return &FeaturesServices{
		CreateFeatures: NewCreateFeaturesService(repository),
		GetFeatures:    NewGetFeaturesService(repository),
		DeleteFeatures: NewDeleteFeaturesService(repository),
		UpdateFeatures: NewUpdateFeaturesService(repository),
		ListFeatures:   NewListFeatures(repository),
	}
}

type Pubsub struct {
	ProcessAggregationInitService          *ProcessAggregationInitService
	ProcessAggregationByDistributorService *ProcessAggregationByDistributorService
}

func NewPubsub(
	publisher event.PublisherCreator,
	topic string,
	inventoryClient clients.Inventory,
	aggregationRepository aggregations.AggregationRepository,
) *Pubsub {
	return &Pubsub{
		ProcessAggregationInitService:          NewProcessAggregationInitService(publisher, inventoryClient, topic),
		ProcessAggregationByDistributorService: NewProcessAggregationByDistributorService(aggregationRepository),
	}
}

type AggregationConfigServices struct {
	GetAggregationConfigsService    *GetAggregationConfigsService
	GetAggregationConfigByIdService *GetAggregationConfigByIdService
	CreateAggregationConfigService  *CreateAggregationConfigService
	UpdateAggregationConfigService  *UpdateAggregationConfigService
	DeleteAggregationConfigService  *DeleteAggregationConfigService
	Location                        *time.Location
}

func NewAggregationConfigServices(aggregationRepository aggregations.AggregationConfigRepository, featuresRepository aggregations.AggregationsFeaturesRepository, schedulerCreator scheduler.ClientCreator, topic string, loc *time.Location) *AggregationConfigServices {
	return &AggregationConfigServices{
		GetAggregationConfigsService:    NewGetAggregationConfigService(aggregationRepository),
		GetAggregationConfigByIdService: NewGetAggregationConfigByIdService(aggregationRepository),
		CreateAggregationConfigService:  NewCreateAggregationConfigService(aggregationRepository, featuresRepository, schedulerCreator, topic, loc),
		UpdateAggregationConfigService:  NewUpdateAggregationConfigService(aggregationRepository, featuresRepository, schedulerCreator, topic, loc),
		DeleteAggregationConfigService:  NewDeleteAggregationConfigService(aggregationRepository, schedulerCreator, topic),
		Location:                        loc,
	}
}

type ServicePoint struct {
	Cups         string `bson:"cups"`
	ServicePoint string `bson:"serial_number"`
}

type AggregationsService struct {
	GetAggregations *GetAggregations
	GetAggregation  *GetAggregation
	Location        *time.Location
}

func NewAggregationsService(repository aggregations.AggregationMongoRepository, location *time.Location) *AggregationsService {
	return &AggregationsService{
		GetAggregations: NewGetAggregationsDashboard(repository, location),
		GetAggregation:  NewGetAggregation(repository, location),
		Location:        location,
	}
}
