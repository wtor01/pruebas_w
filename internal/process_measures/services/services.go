package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"time"
)

type ProcessMeasuresStatsDashboardServices struct {
	GetDashboardMeasuresStats     *GetDashboardMeasuresStats
	GetDashboardMeasuresStatsCups *GetDashboardMeasuresStatsCups
}

func NewProcessMeasuresStatsDashboardServices(repository process_measures.ProcessMeasureDashboardStatsRepository) *ProcessMeasuresStatsDashboardServices {
	return &ProcessMeasuresStatsDashboardServices{
		GetDashboardMeasuresStats:     NewGetDashboardMeasuresStats(repository),
		GetDashboardMeasuresStatsCups: NewGetDashboardMeasuresStatsCups(repository),
	}
}

type DashboardServices struct {
	Location                                   *time.Location
	GetDashboard                               *GetDashboardMeasures
	ListDashboardCups                          *ListDashboardCupsService
	SearchServicePointProcessMeasuresDashboard *SearchServicePointProcessMeasuresDashboard
	SearchServicePointDashboardCurves          *SearchServicePointDashboardCurves
}

func NewDashboardServices(processRepository process_measures.Repository, inventoryRepository measures.InventoryRepository, calendarRepository measures.CalendarPeriodRepository, location *time.Location) *DashboardServices {
	return &DashboardServices{
		Location:          location,
		GetDashboard:      NewGetDashboardMeasures(processRepository, inventoryRepository, location),
		ListDashboardCups: NewListDashboardCupsService(inventoryRepository, processRepository, location),
		SearchServicePointProcessMeasuresDashboard: NewSearchServicePointProcessMeasuresDashboard(processRepository, inventoryRepository, calendarRepository, location),
		SearchServicePointDashboardCurves:          NewSearchServicePointDashboardCurves(processRepository, location),
	}
}

type SchedulerServices struct {
	GetSchedulerById *GetSchedulerById
	CreateScheduler  *CreateScheduler
	DeleteScheduler  *DeleteScheduler
	ListScheduler    *ListScheduler
	UpdateScheduler  *UpdateScheduler
	SearchScheduler  *SearchScheduler
}

func NewSchedulerServices(repository process_measures.SchedulerRepository, schedulerClientCreator scheduler.ClientCreator, topic string) *SchedulerServices {
	return &SchedulerServices{
		GetSchedulerById: NewGetSchedulerByIdService(repository),
		CreateScheduler:  NewCreateScheduler(repository, schedulerClientCreator, topic),
		DeleteScheduler:  NewDeleteScheduler(repository, schedulerClientCreator),
		ListScheduler:    NewListSchedulerService(repository),
		UpdateScheduler:  NewUpdateScheduler(repository, schedulerClientCreator, topic),
		SearchScheduler:  NewSearchSchedulerService(repository),
	}
}

type ServicesPubsub struct {
	PublishDistributorService  *PublishDistributorService
	PublishServicePointService *PublishServicePointService
	ProcessMonthlyClosure      *ProcessMonthlyClosure
	ProcessCurve               *ProcessCurve
	ProcessDailyClosure        *ProcessDailyClosure
}

func NewServicesPubsub(
	publisher event.PublisherCreator,
	topic string,
	topicBilling string,
	repositoryScheduler process_measures.SchedulerRepository,
	inventoryClient clients.Inventory,
	inventoryRepository measures.InventoryRepository,
	grossMeasureRepository gross_measures.GrossMeasureRepository,
	processMeasureRepository process_measures.ProcessedMeasureRepository,
	calendarPeriodRepository measures.CalendarPeriodRepository,
	seasonsRepository seasons.RepositorySeasons,
	location *time.Location,
	validationClient clients.Validation,
	calendar calendar.RepositoryCalendar,
	masterTablesClient clients.MasterTables,
	validationRepository validations.ValidationMongoRepository,
) *ServicesPubsub {

	publishDistributorService := NewPublishDistributorService(
		publisher,
		NewSearchSchedulerService(repositoryScheduler),
		inventoryClient,
		topic,
	)
	publishServicePointService := NewServicePointService(
		publisher,
		inventoryRepository,
		topic,
	)
	processMonthlyClosure := NewProcessMonthlyClosure(
		grossMeasureRepository,
		processMeasureRepository,
		calendar,
		validationClient,
		calendarPeriodRepository,
		location,
		masterTablesClient,
		validationRepository,
		publisher,
		topicBilling,
	)
	processCurve := NewProcessCurve(
		grossMeasureRepository,
		processMeasureRepository,
		calendarPeriodRepository,
		seasonsRepository,
		location,
		validationClient,
		masterTablesClient,
		validationRepository,
	)
	processDailyClosure := NewProcessDailyClosure(
		grossMeasureRepository,
		processMeasureRepository,
		validationClient,
		calendarPeriodRepository,
		location,
		masterTablesClient,
		validationRepository,
	)
	return &ServicesPubsub{
		PublishDistributorService:  publishDistributorService,
		PublishServicePointService: publishServicePointService,
		ProcessMonthlyClosure:      processMonthlyClosure,
		ProcessCurve:               processCurve,
		ProcessDailyClosure:        processDailyClosure,
	}
}

type ReprocessingServicesPubSub struct {
	ReprocessingPublishDistributorService  *ReprocessingPublishDistributorService
	PublishReprocessingServicePointService *PublishReprocessingServicePointService
	PublishReprocessingMeterService        *PublishReprocessingMeterService
}

func NewReprocessingServicesPubSub(
	publisher event.PublisherCreator,
	topic string,
	repositoryScheduler process_measures.SchedulerRepository,
	inventoryClient clients.Inventory,
	inventoryRepository measures.InventoryRepository,
	grossMeasureRepository gross_measures.GrossMeasureRepository,
	reprocessingDateRepository measures.ReprocessingDateRepository) *ReprocessingServicesPubSub {
	publishReprocessingServicePointService := NewReprocessingServicePointService(
		publisher,
		topic,
		grossMeasureRepository,
		100,
	)
	publishReprocessingMeterService := NewReprocessingMeterService(
		publisher,
		inventoryRepository,
		topic,
	)
	reprocessingPublishDistributor := NewReprocessingPublishDistributorService(
		publisher,
		NewSearchSchedulerService(repositoryScheduler),
		inventoryClient,
		topic,
		reprocessingDateRepository,
		time.Now,
	)

	return &ReprocessingServicesPubSub{
		PublishReprocessingServicePointService: publishReprocessingServicePointService,
		ReprocessingPublishDistributorService:  reprocessingPublishDistributor,
		PublishReprocessingMeterService:        publishReprocessingMeterService,
	}
}
