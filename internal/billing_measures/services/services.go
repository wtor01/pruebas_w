package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"time"
)

type SchedulerServices struct {
	ListScheduler    *ListScheduler
	GetSchedulerById *GetSchedulerById
	DeleteScheduler  *DeleteScheduler
	CreateScheduler  *CreateScheduler
	UpdateScheduler  *UpdateScheduler
}

func NewSchedulerServices(repository billing_measures.BillingSchedulerRepository, schedulerClientCreator scheduler.ClientCreator, topic string) *SchedulerServices {
	return &SchedulerServices{
		ListScheduler:    NewListSchedulerService(repository),
		GetSchedulerById: NewGetSchedulerByIdService(repository),
		DeleteScheduler:  NewDeleteScheduler(repository, schedulerClientCreator),
		CreateScheduler:  NewCreateScheduler(repository, schedulerClientCreator, topic),
		UpdateScheduler:  NewUpdateScheduler(repository, schedulerClientCreator, topic),
	}
}

type ServicesPubsub struct {
	PublishDistributorService  *PublishDistributorService
	PublishServicePointService *PublishServicePointService
	ProcessDcMeasureByCup      *ProcessMvhByCup
	ProcessSelfConsumption     *ProcessSelfConsumption
}

func NewServicesPubsub(
	publisher event.PublisherCreator,
	topic string,
	billingSchedulerRepository billing_measures.BillingSchedulerRepository,
	billingMeasureRepository billing_measures.BillingMeasureRepository,
	selfConsumptionRepository billing_measures.SelfConsumptionRepository,
	billingSelfConsumptionRepository billing_measures.BillingSelfConsumptionRepository,
	processedMeasureRepository process_measures.ProcessedMeasureRepository,
	inventoryClient clients.Inventory,
	inventoryRepository measures.InventoryRepository,
	repositoryProfiles billing_measures.ConsumProfileRepository,
	calendarPeriodRepository measures.CalendarPeriodRepository,
	location *time.Location,
	masterTablesClient clients.MasterTables,
	consumptionCoefficientRepository billing_measures.ConsumCoefficientRepository,
) *ServicesPubsub {

	publishDistributorService := NewPublishDistributorService(publisher, NewSearchSchedulerService(billingSchedulerRepository), inventoryClient, topic)
	publishServicePointService := NewServicePointService(publisher, inventoryRepository, topic)
	processTgMeasureByCup := NewProcessMvhByCup(
		billingMeasureRepository,
		processedMeasureRepository,
		inventoryClient,
		repositoryProfiles,
		calendarPeriodRepository,
		location,
		masterTablesClient,
		publisher,
		topic,
	)

	return &ServicesPubsub{
		PublishDistributorService:  publishDistributorService,
		PublishServicePointService: publishServicePointService,
		ProcessDcMeasureByCup:      processTgMeasureByCup,
		ProcessSelfConsumption: NewProcessSelfConsumption(
			billingMeasureRepository,
			selfConsumptionRepository,
			billingSelfConsumptionRepository,
			location,
			consumptionCoefficientRepository,
		),
	}
}

type BillingMeasuresDashboardService struct {
	Location                    *time.Location
	TaxMeasuresByCups           *TaxMeasuresByCups
	SearchFiscalBillingMeasures *SearchFiscalBillingMeasures
	DashboardSummaryService     *DashboardSummaryService
	ExecuteMvh                  *ServiceExecuteMvh
	GetClosureResumeDashboard   *GetClosureResumeDashboard
}

func NewBillingMeasuresDashboardService(billingMeasureRepository billing_measures.BillingMeasuresDashboardRepository, location *time.Location, inventoryRepository measures.InventoryRepository, publisher event.PublisherCreator, topic string) *BillingMeasuresDashboardService {
	return &BillingMeasuresDashboardService{
		SearchFiscalBillingMeasures: NewSearchFiscalBillingMeasuresDashboard(billingMeasureRepository, location),
		TaxMeasuresByCups:           NewTaxMeasuresByCups(billingMeasureRepository),
		ExecuteMvh:                  NewServiceExecuteMvh(inventoryRepository, publisher, topic),
		GetClosureResumeDashboard:   NewGetClosureResumeDashboard(billingMeasureRepository),
		DashboardSummaryService:     NewDashboardSummaryService(billingMeasureRepository),
		Location:                    location,
	}
}
