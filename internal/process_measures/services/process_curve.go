package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"time"
)

type ProcessCurve struct {
	ProcessMeasureBase
	repository             process_measures.ProcessedMeasureRepository
	grossMeasureRepository gross_measures.GrossMeasureRepository
	seasonsRepository      seasons.RepositorySeasons
	validationRepository   validations.ValidationMongoRepository
}

func NewProcessCurve(
	grossMeasureRepository gross_measures.GrossMeasureRepository,
	repository process_measures.ProcessedMeasureRepository,
	calendarPeriodRepository measures.CalendarPeriodRepository,
	seasonsRepository seasons.RepositorySeasons,
	Location *time.Location,
	validationClient clients.Validation,
	masterTablesClient clients.MasterTables,
	validationRepository validations.ValidationMongoRepository,

) *ProcessCurve {
	return &ProcessCurve{
		ProcessMeasureBase: NewProcessMeasureBase(
			calendarPeriodRepository,
			Location,
			validationClient,
			masterTablesClient,
		),
		grossMeasureRepository: grossMeasureRepository,
		repository:             repository,
		seasonsRepository:      seasonsRepository,
		validationRepository:   validationRepository,
	}

}

func (svc ProcessCurve) Handle(ctx context.Context, dto measures.ProcessMeasurePayload) error {
	ms, err := svc.grossMeasureRepository.ListDailyCurveMeasures(ctx, gross_measures.QueryListForProcessCurve{
		SerialNumber: dto.MeterConfig.Meter.SerialNumber,
		Date:         dto.Date,
		CurveType:    dto.CurveType,
	})
	if err != nil {
		return err
	}

	measureByDate := svc.GenerateMeasuresCurveEmpty(ctx, dto)

	tariff, err := svc.GetTariff(ctx, dto)
	calendarPeriodDay, err := svc.GetCalendarPeriod(ctx, dto, tariff)

	if err != nil {
		return err
	}

	validators, err := svc.GetValidations(ctx, dto, svc.validationRepository)

	if err != nil {
		return err
	}

	dayType, err := svc.seasonsRepository.GetDayTypeByMonth(ctx, int(dto.Date.Month()), calendarPeriodDay.IsFestiveDay())

	if err != nil {
		return err
	}
	measureByDate = svc.FillMeasuresCurve(ctx, dto, measureByDate, ms, calendarPeriodDay, dayType)
	measureByDate = svc.FillEmptyMeasuresCurve(ctx, dto, measureByDate, calendarPeriodDay, dayType)
	processMeasuresCurve := process_measures.ListProcessedLoadCurve(utils.MapToSlice[process_measures.ProcessedLoadCurve](measureByDate))

	svc.ValidateMeasures(processMeasuresCurve, validators)

	err = svc.repository.SaveAllProcessedLoadCurve(ctx, processMeasuresCurve)

	if err != nil {
		return err
	}

	qualityCode := processMeasuresCurve.GetQualityCodeLoadCurve()

	var m process_measures.ProcessedLoadCurve

	if len(processMeasuresCurve) != 0 {
		m = processMeasuresCurve[0]
	}

	daily := process_measures.ProcessedDailyLoadCurve{
		DistributorID:            dto.MeterConfig.DistributorID,
		DistributorCode:          dto.MeterConfig.DistributorCode,
		CUPS:                     dto.MeterConfig.Cups(),
		EndDate:                  dto.Date.UTC().AddDate(0, 0, 1),
		MeterSerialNumber:        dto.MeterConfig.SerialNumber(),
		GenerationDate:           svc.generatorDate().UTC(),
		ReadingDate:              m.ReadingDate,
		RegisterType:             m.RegisterType,
		CurveType:                m.CurveType,
		ServiceType:              dto.MeterConfig.ServiceType(),
		PointType:                dto.MeterConfig.PointType(),
		MeterType:                dto.MeterConfig.MeterType(),
		ValidationStatus:         measures.Valid,
		ValidationStatusIsManual: false,
		InvalidationCodes:        nil,
		QualityCode:              qualityCode,
	}

	daily.GenerateID()

	err = svc.repository.SaveProcessedDailyLoadCurve(ctx, daily)

	return err
}

func (svc ProcessCurve) GenerateMeasuresCurveEmpty(ctx context.Context, dto measures.ProcessMeasurePayload) map[string]process_measures.ProcessedLoadCurve {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GenerateMeasuresCurveEmpty")
	defer span.End()

	delta := time.Hour
	if dto.CurveType == measures.QuarterMeasureCurveReadingType {
		delta = time.Minute * 15
	}
	startDate := dto.Date.UTC().Add(delta)
	endDate := dto.Date.UTC().AddDate(0, 0, 1)
	measureByDate := make(map[string]process_measures.ProcessedLoadCurve)

	for startDate.Before(endDate) || startDate.Equal(endDate) {
		measureByDate[startDate.Format("2006-01-02 15:04")] = process_measures.ProcessedLoadCurve{}
		startDate = startDate.Add(delta)
	}

	return measureByDate
}

func (svc ProcessCurve) FillMeasuresCurve(ctx context.Context, dto measures.ProcessMeasurePayload, measureByDate map[string]process_measures.ProcessedLoadCurve, ms []gross_measures.MeasureCurveWrite, calendarPeriodDay measures.CalendarPeriod, dayType seasons.DayTypes) map[string]process_measures.ProcessedLoadCurve {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "FillMeasuresCurve")
	defer span.End()

	magnitudes := dto.MeterConfig.GetMagnitudesActive()

	for _, m := range ms {
		key := m.EndDate.Format("2006-01-02 15:04")
		if v, ok := measureByDate[key]; ok && !v.GenerationDate.IsZero() {
			continue
		}
		processedLoadCurve := process_measures.NewProcessedLoadCurve(
			dto.MeterConfig,
			m.EndDate,
			dto.CurveType,
		)

		processedLoadCurve.GenerationDate = svc.generatorDate().UTC()
		processedLoadCurve.ReadingDate = m.ReadingDate
		processedLoadCurve.Origin = m.Origin
		processedLoadCurve.AI = m.AI
		processedLoadCurve.AE = m.AE
		processedLoadCurve.R1 = m.R1
		processedLoadCurve.R2 = m.R2
		processedLoadCurve.R3 = m.R3
		processedLoadCurve.R4 = m.R4
		processedLoadCurve.Period = calendarPeriodDay.GetHourPeriod(m.EndDate, dto.CurveType, dayType.IsFestive, svc.Location)
		processedLoadCurve.Magnitudes = magnitudes
		processedLoadCurve.DayTypeId = dayType.ID
		processedLoadCurve.SeasonId = dayType.SeasonsId

		measureByDate[key] = *processedLoadCurve
	}

	return measureByDate
}

func (svc ProcessCurve) FillEmptyMeasuresCurve(ctx context.Context, dto measures.ProcessMeasurePayload, measureByDate map[string]process_measures.ProcessedLoadCurve, calendarPeriodDay measures.CalendarPeriod, dayType seasons.DayTypes) map[string]process_measures.ProcessedLoadCurve {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "FillEmptyMeasuresCurve")
	defer span.End()

	magnitudes := dto.MeterConfig.GetMagnitudesActive()

	for k, v := range measureByDate {
		if !v.GenerationDate.IsZero() {
			continue
		}
		endDateFilled, err := time.Parse("2006-01-02 15:04", k)
		if err != nil {
			continue
		}

		processedLoadCurve := process_measures.NewProcessedLoadCurve(
			dto.MeterConfig,
			endDateFilled.In(time.UTC),
			dto.CurveType,
		)
		processedLoadCurve.GenerationDate = svc.generatorDate().UTC()
		processedLoadCurve.Origin = measures.Filled
		processedLoadCurve.Period = calendarPeriodDay.GetHourPeriod(endDateFilled, dto.CurveType, dayType.IsFestive, svc.Location)
		processedLoadCurve.Magnitudes = magnitudes
		processedLoadCurve.DayTypeId = dayType.ID
		processedLoadCurve.SeasonId = dayType.SeasonsId

		measureByDate[k] = *processedLoadCurve
	}

	return measureByDate
}
