package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

type ProcessMonthlyClosure struct {
	ProcessMeasureBase
	grossMeasureRepository   gross_measures.GrossMeasureRepository
	processMeasureRepository process_measures.ProcessedMeasureRepository
	calendar                 calendar.RepositoryCalendar
	validationRepository     validations.ValidationMongoRepository
	publish                  event.PublisherCreator
	topic                    string
}

func NewProcessMonthlyClosure(
	grossMeasureRepository gross_measures.GrossMeasureRepository,
	processMeasureRepository process_measures.ProcessedMeasureRepository,
	calendar calendar.RepositoryCalendar,
	validationClient clients.Validation,
	calendarPeriodRepository measures.CalendarPeriodRepository,
	Location *time.Location,
	masterTablesClient clients.MasterTables,
	validationRepository validations.ValidationMongoRepository,
	publish event.PublisherCreator,
	topic string,
) *ProcessMonthlyClosure {
	return &ProcessMonthlyClosure{
		ProcessMeasureBase: NewProcessMeasureBase(
			calendarPeriodRepository,
			Location,
			validationClient,
			masterTablesClient,
		),
		grossMeasureRepository:   grossMeasureRepository,
		processMeasureRepository: processMeasureRepository,
		calendar:                 calendar,
		validationRepository:     validationRepository,
		publish:                  publish,
		topic:                    topic,
	}

}

func (svc ProcessMonthlyClosure) Handle(ctx context.Context, dto measures.ProcessMeasurePayload) error {
	ms, err := svc.grossMeasureRepository.ListDailyCloseMeasures(ctx, gross_measures.QueryListForProcessClose{
		ReadingType:  measures.BillingClosure,
		SerialNumber: dto.MeterConfig.SerialNumber(),
		Date:         dto.Date,
	})

	if err != nil {
		return err
	}

	if len(ms) == 0 {
		return nil
	}

	bestMeasureByDate := make(map[string]struct {
		increment gross_measures.MeasureCloseWrite
		absolute  gross_measures.MeasureCloseWrite
	})

	//Obtener la medida procesada del mes anterior
	processedMeasureLM, _ := svc.processMeasureRepository.GetMonthlyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date.AddDate(0, -1, 1),
	})

	//Escoger las mejores medidas de incremento y absolutas para ese meter Id por hora
	for _, measure := range ms {
		dateString := measure.EndDate.Format("2006-01-02 15")

		info, _ := bestMeasureByDate[dateString]

		switch measure.Type {
		case measures.Absolute:
			if info.absolute.Id == "" {
				info.absolute = measure
			}
		case measures.Incremental:
			if info.increment.Id == "" {
				info.increment = measure
			}
		}
		bestMeasureByDate[dateString] = info
	}

	tariff, err := svc.GetTariff(ctx, dto)

	if err != nil {
		return err
	}

	grouperr, ctx := errgroup.WithContext(ctx)

	for _, measuresInfo := range bestMeasureByDate {
		measuresInfo := measuresInfo
		grouperr.Go(func() error {
			incrementMeasure := measuresInfo.increment

			absoluteMeasure := svc.fillAbsoluteMeasureIfNeeded(measuresInfo.absolute, incrementMeasure)
			calendarPeriods := svc.GetProcessedMonthlyClosureCalendar(absoluteMeasure, incrementMeasure)

			processMeasureMonthlyClosure := process_measures.NewProcessedMonthlyClosure(dto.MeterConfig, absoluteMeasure.EndDate)
			processMeasureMonthlyClosure.GenerationDate = svc.generatorDate().UTC()
			processMeasureMonthlyClosure.ReadingDate = absoluteMeasure.ReadingDate
			processMeasureMonthlyClosure.StartDate = absoluteMeasure.StartDate
			processMeasureMonthlyClosure.Origin = absoluteMeasure.Origin
			processMeasureMonthlyClosure.ContractNumber = absoluteMeasure.Contract
			processMeasureMonthlyClosure.CalendarPeriods = calendarPeriods
			processMeasureMonthlyClosure.Coefficient = tariff.Coef

			calendarPeriodDay, errGo := svc.calendarPeriodRepository.GetCalendarPeriod(ctx, measures.SearchCalendarPeriod{
				Day:          dto.Date,
				GeographicID: tariff.GeographicId,
				CalendarCode: tariff.CalendarId,
				Location:     svc.Location,
			})

			if errGo != nil {
				return errGo
			}

			periods := calendarPeriodDay.GetAllPeriods()

			processMeasureMonthlyClosure.Periods = periods
			processMeasureMonthlyClosure.SetLastDailyClose(processedMeasureLM)

			validators, errGo := svc.GetValidations(ctx, dto, svc.validationRepository)

			svc.ValidateMeasure(processMeasureMonthlyClosure, validators)

			errGo = svc.processMeasureRepository.SaveMonthlyClosure(ctx, *processMeasureMonthlyClosure)

			if errGo != nil {
				return errGo
			}

			endDateWithLocation := processMeasureMonthlyClosure.EndDate.In(svc.Location)

			if endDateWithLocation.Day() == 1 || endDateWithLocation.Hour() != 0 {
				return nil
			}

			e := billing_measures.NewProcessMvhEvent(endDateWithLocation, dto.MeterConfig)

			errGo = event.PublishEvent(ctx, svc.topic, svc.publish, e)

			return errGo
		})

	}

	err = grouperr.Wait()

	return err
}

func (svc ProcessMonthlyClosure) GetProcessedMonthlyClosureCalendar(bestAbsoluteMeasure, bestIncrementClose gross_measures.MeasureCloseWrite) process_measures.ProcessedMonthlyClosureCalendar {

	var calendarPeriods process_measures.ProcessedMonthlyClosureCalendar

	for _, p := range append(measures.ValidPeriodsCurve, measures.P0) {
		switch p {
		case measures.P0:
			calendarPeriods.P0 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P1:
			calendarPeriods.P1 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P2:
			calendarPeriods.P2 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P3:
			calendarPeriods.P3 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P4:
			calendarPeriods.P4 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P5:
			calendarPeriods.P5 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P6:
			calendarPeriods.P6 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		}
	}

	for _, period := range bestAbsoluteMeasure.Periods {

		var processedMonthlyClosurePeriod *process_measures.ProcessedMonthlyClosurePeriod

		switch period.Period {
		case measures.P0:
			processedMonthlyClosurePeriod = calendarPeriods.P0
		case measures.P1:
			processedMonthlyClosurePeriod = calendarPeriods.P1
		case measures.P2:
			processedMonthlyClosurePeriod = calendarPeriods.P2
		case measures.P3:
			processedMonthlyClosurePeriod = calendarPeriods.P3
		case measures.P4:
			processedMonthlyClosurePeriod = calendarPeriods.P4
		case measures.P5:
			processedMonthlyClosurePeriod = calendarPeriods.P5
		case measures.P6:
			processedMonthlyClosurePeriod = calendarPeriods.P6
		}

		if processedMonthlyClosurePeriod == nil {
			continue
		}
		processedMonthlyClosurePeriod.Values = measures.Values{
			AI: period.AI,
			AE: period.AE,
			R1: period.R1,
			R2: period.R2,
			R3: period.R3,
			R4: period.R4,
		}
		processedMonthlyClosurePeriod.ValidationStatus = measures.Valid
		processedMonthlyClosurePeriod.Filled = false
	}

	for _, period := range bestIncrementClose.Periods {
		var calendarPeriod *process_measures.ProcessedMonthlyClosurePeriod
		switch period.Period {
		case measures.P0:
			calendarPeriod = calendarPeriods.P0
		case measures.P1:
			calendarPeriod = calendarPeriods.P1
		case measures.P2:
			calendarPeriod = calendarPeriods.P2
		case measures.P3:
			calendarPeriod = calendarPeriods.P3
		case measures.P4:
			calendarPeriod = calendarPeriods.P4
		case measures.P5:
			calendarPeriod = calendarPeriods.P5
		case measures.P6:
			calendarPeriod = calendarPeriods.P6
		}
		if calendarPeriod == nil {
			continue
		}
		calendarPeriod.AIi = period.AI
		calendarPeriod.AEi = period.AE
		calendarPeriod.R1i = period.R1
		calendarPeriod.R2i = period.R2
		calendarPeriod.R3i = period.R3
		calendarPeriod.R4i = period.R4
	}

	return calendarPeriods
}

func (svc ProcessMonthlyClosure) fillAbsoluteMeasureIfNeeded(
	absoluteMeasure, incrementMeasure gross_measures.MeasureCloseWrite,
) gross_measures.MeasureCloseWrite {
	if absoluteMeasure.Id != "" {
		return absoluteMeasure
	}
	absoluteMeasure.EndDate = incrementMeasure.EndDate
	absoluteMeasure.ReadingDate = incrementMeasure.ReadingDate
	absoluteMeasure.StartDate = incrementMeasure.StartDate
	absoluteMeasure.Origin = incrementMeasure.Origin
	absoluteMeasure.Contract = incrementMeasure.Contract
	absoluteMeasure.Periods = make([]gross_measures.MeasureClosePeriod, 0, cap(incrementMeasure.Periods))
	for _, period := range incrementMeasure.Periods {
		absoluteMeasure.Periods = append(absoluteMeasure.Periods, gross_measures.MeasureClosePeriod{
			Period: period.Period,
		})
	}

	return absoluteMeasure
}
