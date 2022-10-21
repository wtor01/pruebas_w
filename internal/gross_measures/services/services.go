package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	"time"
)

type Services struct {
	InsertMeasureCurve *InsertMeasureCurveService
	InsertMeasureClose *InsertMeasureCloseService
	ParseFilesService  ParseFilesService
}

func NewServices(
	measureRepository gross_measures.GrossMeasureRepository,
	publisher event.PublisherCreator,
	storage storage.StorageCreator,
	cnf config.Config,
	clientInventory clients.Inventory,
	clientValidation clients.Validation,
) *Services {

	insertMeasureCurve := NewInsertMeasureCurveService(measureRepository, clientValidation, clientInventory, publisher, cnf.ProcessedMeasureTopic)
	insertMeasureClose := NewInsertMeasureCloseService(measureRepository, clientValidation, clientInventory, publisher, cnf.ProcessedMeasureTopic)
	parseFilesService := NewParseFilesService(publisher, cnf.TopicMeasures, storage, cnf.StorageMeasuresFail, cnf.StorageMeasuresSuccess, cnf.LocalLocation)

	return &Services{
		InsertMeasureCurve: insertMeasureCurve,
		InsertMeasureClose: insertMeasureClose,
		ParseFilesService:  parseFilesService,
	}
}

type DashboardServices struct {
	DashboardMeasureSupplyPointService *DashboardMeasureSupplyPointService
	DashboardSupplyPointCurvesService  *DashboardSupplyPointCurvesService
}

func NewDashboardServices(grossRepository gross_measures.GrossMeasureRepository, inventoryRepository measures.InventoryRepository, calendarPeriodRepository measures.CalendarPeriodRepository, masterTablesClient clients.MasterTables, loc *time.Location) *DashboardServices {
	return &DashboardServices{
		DashboardMeasureSupplyPointService: NewDashboardMeasureSupplyPointService(
			grossRepository,
			inventoryRepository,
			calendarPeriodRepository,
			masterTablesClient,
			loc),
		DashboardSupplyPointCurvesService: NewDashboardSupplyPointCurvesService(
			grossRepository,
			inventoryRepository,
			loc),
	}
}

type GrossMeasuresStatsDashboardServices struct {
	GetDashboardMeasuresStats              *GetDashboardMeasuresStats
	ListDashboardMeasuresStatsSerialNumber *ListDashboardMeasuresStatsSerialNumber
}

func NewGrossMeasuresStatsDashboardServices(repository gross_measures.GrossMeasuresDashboardStatsRepository) *GrossMeasuresStatsDashboardServices {
	return &GrossMeasuresStatsDashboardServices{
		GetDashboardMeasuresStats:              NewGetDashboardMeasuresStats(repository),
		ListDashboardMeasuresStatsSerialNumber: NewListDashboardMeasuresStatsSerialNumber(repository),
	}
}
