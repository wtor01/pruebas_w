package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"time"
)

type CreateClosure struct {
	ProcessMeasureBase
	Location            *time.Location
	closureRepository   process_measures.ProcessMeasureClosureRepository
	inventoryRepository measures.InventoryRepository
	generatorDate       func() time.Time
	masterTablesClient  clients.MasterTables
}

func NewCreateClosure(closureRepository process_measures.ProcessMeasureClosureRepository,
	inventoryRepository measures.InventoryRepository,
	masterTablesClient clients.MasterTables,
	generatorDate func() time.Time,
) *CreateClosure {
	return &CreateClosure{closureRepository: closureRepository, inventoryRepository: inventoryRepository,
		masterTablesClient: masterTablesClient,
		generatorDate:      generatorDate,
	}
}

type CreateClosureDto struct {
	Monthly process_measures.ProcessedMonthlyClosure
}

func (c CreateClosure) Handler(ctx context.Context, dto CreateClosureDto) error {

	monthly := dto.Monthly

	monthly.GenerationDate = c.generatorDate()

	resultMeterConfig, err := c.inventoryRepository.GetMeterConfigByCupsAPI(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        monthly.CUPS,
		Distributor: monthly.DistributorID,
		Time:        monthly.EndDate,
	})

	if err != nil {
		return err
	}

	monthly.ContractNumber = string(resultMeterConfig.PriorityContract)
	processMeasureMonthlyClosure := process_measures.NewProcessedMonthlyClosure(resultMeterConfig, time.Time{})
	processMeasureMonthlyClosure.ContractNumber = string(resultMeterConfig.PriorityContract)
	processMeasureMonthlyClosure.StartDate = monthly.StartDate
	processMeasureMonthlyClosure.EndDate = monthly.EndDate
	processMeasureMonthlyClosure.GenerationDate = c.generatorDate()
	processMeasureMonthlyClosure.ReadingDate = c.generatorDate()

	tariff, err := c.masterTablesClient.GetTariff(ctx, clients.GetTariffDto{
		ID: resultMeterConfig.TariffID(),
	})
	if err != nil {
		return err
	}

	monthly.Coefficient = tariff.Coef
	processMeasureMonthlyClosure.Origin = measures.Manual
	processMeasureMonthlyClosure.CalendarPeriods = monthly.CalendarPeriods
	processMeasureMonthlyClosure.Coefficient = tariff.Coef
	processMeasureMonthlyClosure.Id = monthly.Id

	err = c.closureRepository.CreateClosure(ctx, *processMeasureMonthlyClosure)

	return err

}
