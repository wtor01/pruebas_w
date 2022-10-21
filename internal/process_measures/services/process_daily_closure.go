package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
	"errors"
	"time"
)

type ProcessDailyClosure struct {
	ProcessMeasureBase
	grossMeasureRepository   gross_measures.GrossMeasureRepository
	processMeasureRepository process_measures.ProcessedMeasureRepository
	validationRepository     validations.ValidationMongoRepository
}

func NewProcessDailyClosure(
	grossMeasureRepository gross_measures.GrossMeasureRepository,
	processMeasureRepository process_measures.ProcessedMeasureRepository,
	validationClient clients.Validation,
	calendarPeriodRepository measures.CalendarPeriodRepository,
	Location *time.Location,
	masterTablesClient clients.MasterTables,
	validationRepository validations.ValidationMongoRepository,

) *ProcessDailyClosure {
	return &ProcessDailyClosure{
		ProcessMeasureBase: NewProcessMeasureBase(
			calendarPeriodRepository,
			Location,
			validationClient,
			masterTablesClient,
		),
		grossMeasureRepository:   grossMeasureRepository,
		processMeasureRepository: processMeasureRepository,
		validationRepository:     validationRepository,
	}
}

func (svc ProcessDailyClosure) Handle(ctx context.Context, dto measures.ProcessMeasurePayload) error {
	ms, err := svc.grossMeasureRepository.ListDailyCloseMeasures(ctx, gross_measures.QueryListForProcessClose{
		ReadingType:  measures.DailyClosure,
		SerialNumber: dto.MeterConfig.SerialNumber(),
		Date:         dto.Date,
	})

	if err != nil {
		return err
	}

	if len(ms) == 0 {
		return nil
	}

	bestDailyMeasure := ms[0]

	if bestDailyMeasure.Type == "" {
		return errors.New("invalid measure")
	}

	processedMeasureLM, _ := svc.processMeasureRepository.GetProcessedDailyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date.AddDate(0, 0, 1),
	})

	tariff, err := svc.GetTariff(ctx, dto)

	if err != nil {
		return err
	}

	calendarPeriods := svc.GetProcessedDailyClosureCalendar(bestDailyMeasure)

	measureProcessed := process_measures.NewProcessedDailyClosure(dto.MeterConfig, bestDailyMeasure.EndDate)

	measureProcessed.DistributorCode = bestDailyMeasure.DistributorCDOS
	measureProcessed.GenerationDate = svc.generatorDate().UTC()
	measureProcessed.ReadingDate = bestDailyMeasure.ReadingDate
	measureProcessed.Origin = bestDailyMeasure.Origin
	measureProcessed.ContractNumber = bestDailyMeasure.Contract
	measureProcessed.CalendarPeriods = calendarPeriods
	measureProcessed.ValidationStatus = bestDailyMeasure.Status
	measureProcessed.Coefficient = tariff.Coef

	measureProcessed.SetLastDailyClose(processedMeasureLM)

	if err != nil {
		return err
	}

	validators, err := svc.GetValidations(ctx, dto, svc.validationRepository)

	svc.ValidateMeasure(measureProcessed, validators)

	err = svc.processMeasureRepository.SaveDailyClosure(ctx, *measureProcessed)

	return err
}

func (svc ProcessDailyClosure) GetProcessedDailyClosureCalendar(bestDailyMeasure gross_measures.MeasureCloseWrite) process_measures.ProcessedDailyClosureCalendar {

	var calendarPeriods process_measures.ProcessedDailyClosureCalendar
	calendarPeriods.P0 = nil
	calendarPeriods.P1 = nil
	calendarPeriods.P2 = nil
	calendarPeriods.P3 = nil
	calendarPeriods.P4 = nil
	calendarPeriods.P5 = nil
	calendarPeriods.P6 = nil

	for _, p := range append(measures.ValidPeriodsCurve, measures.P0) {
		switch p {
		case measures.P0:
			calendarPeriods.P0 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		case measures.P1:
			calendarPeriods.P1 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		case measures.P2:
			calendarPeriods.P2 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		case measures.P3:
			calendarPeriods.P3 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		case measures.P4:
			calendarPeriods.P4 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		case measures.P5:
			calendarPeriods.P5 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		case measures.P6:
			calendarPeriods.P6 = &process_measures.ProcessedDailyClosurePeriod{
				Filled:           true,
				ValidationStatus: measures.Invalid,
			}
		}
	}

	for _, period := range bestDailyMeasure.Periods {
		var processedDailyClosurePeriod *process_measures.ProcessedDailyClosurePeriod

		switch period.Period {
		case measures.P0:
			processedDailyClosurePeriod = calendarPeriods.P0
		case measures.P1:
			processedDailyClosurePeriod = calendarPeriods.P1
		case measures.P2:
			processedDailyClosurePeriod = calendarPeriods.P2
		case measures.P3:
			processedDailyClosurePeriod = calendarPeriods.P3
		case measures.P4:
			processedDailyClosurePeriod = calendarPeriods.P4
		case measures.P5:
			processedDailyClosurePeriod = calendarPeriods.P5
		case measures.P6:
			processedDailyClosurePeriod = calendarPeriods.P6
		}

		if processedDailyClosurePeriod != nil {
			processedDailyClosurePeriod.Values = measures.Values{
				AI: period.AI,
				AE: period.AE,
				R1: period.R1,
				R2: period.R2,
				R3: period.R3,
				R4: period.R4,
			}
			processedDailyClosurePeriod.ValidationStatus = measures.Valid
			processedDailyClosurePeriod.Filled = false
		}
	}

	return calendarPeriods
}
