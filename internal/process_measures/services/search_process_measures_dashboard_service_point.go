package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sort"
	"time"
)

type SearchServicePointProcessMeasuresDashboardDTO struct {
	DistributorId string
	Cups          string
	StartDate     time.Time
	EndDate       time.Time
}

func NewSearchServicePointProcessMeasuresDashboardDTO(distributorId, cups string, startDate, endDate time.Time) SearchServicePointProcessMeasuresDashboardDTO {
	return SearchServicePointProcessMeasuresDashboardDTO{
		DistributorId: distributorId,
		Cups:          cups,
		StartDate:     startDate,
		EndDate:       endDate,
	}
}

type SearchServicePointProcessMeasuresDashboard struct {
	repositoryProcessMeasures process_measures.ProcessedMeasureRepository
	repositoryInventory       measures.InventoryRepository
	calendarPeriodRepository  measures.CalendarPeriodRepository
	tracer                    trace.Tracer
	Location                  *time.Location
}

func NewSearchServicePointProcessMeasuresDashboard(repositoryProcessMeasures process_measures.ProcessedMeasureRepository, repositoryInventory measures.InventoryRepository, calendarPeriodRepository measures.CalendarPeriodRepository, location *time.Location) *SearchServicePointProcessMeasuresDashboard {
	return &SearchServicePointProcessMeasuresDashboard{repositoryProcessMeasures: repositoryProcessMeasures, repositoryInventory: repositoryInventory, calendarPeriodRepository: calendarPeriodRepository, Location: location, tracer: telemetry.GetTracer()}
}

func (s SearchServicePointProcessMeasuresDashboard) Handler(ctx context.Context, dto SearchServicePointProcessMeasuresDashboardDTO) (process_measures.ServicePointDashboardWithType, error) {

	ctx, span := s.tracer.Start(ctx, "SearchServicePointProcessMeasuresDashboard - Handler")

	defer span.End()

	meterConfig, err := s.repositoryInventory.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		Distributor: dto.DistributorId,
		CUPS:        dto.Cups,
		Time:        dto.StartDate,
	})

	if err != nil {
		return process_measures.ServicePointDashboardWithType{}, err
	}
	calendarPeriodDay, err := s.calendarPeriodRepository.GetCalendarPeriod(ctx, measures.SearchCalendarPeriod{
		Day:          dto.StartDate,
		GeographicID: "ES",
		CalendarCode: meterConfig.CalendarID,
		Location:     s.Location,
	})

	if err != nil {
		return process_measures.ServicePointDashboardWithType{}, err
	}

	periods := calendarPeriodDay.GetAllPeriods()
	allPeriods := append(periods, measures.P0)

	sort.Slice(periods, func(i, j int) bool {
		return periods[i] > periods[j]
	})

	curves, err := s.repositoryProcessMeasures.ProcessedLoadCurveByCups(ctx, process_measures.QueryProcessedLoadCurveByCups{
		CUPS:      dto.Cups,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
		CurveType: measures.HourlyMeasureCurveReadingType,
	})
	if err != nil {
		return process_measures.ServicePointDashboardWithType{}, err
	}

	dailyClosureMeasures, err := s.repositoryProcessMeasures.ProcessedDailyClosureByCups(ctx, process_measures.QueryHistoryDailyClosureByCups{
		CUPS:      dto.Cups,
		StartDate: dto.StartDate.AddDate(0, 0, -1),
		EndDate:   dto.EndDate,
	})
	if err != nil {
		return process_measures.ServicePointDashboardWithType{}, err
	}

	monthlyClosureMeasures, err := s.repositoryProcessMeasures.GetMonthlyClosureMeasuresByCup(ctx, process_measures.QueryMonthlyClosedMeasures{
		CUPS:      dto.Cups,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
	})

	if err != nil {
		return process_measures.ServicePointDashboardWithType{}, err
	}

	magnitudeEnergy := measures.AE
	if meterConfig.ServicePoint.ServiceType == measures.DcServiceType {
		magnitudeEnergy = measures.AI
	}

	magnitudes := meterConfig.GetMagnitudesActive()

	result := make([]process_measures.ServicePointDashboard, 0)
	curveMap := make(map[string]*process_measures.MeasureServicePointDashboard)
	dateServicePointDashboardMap := make(map[string]*process_measures.ServicePointDashboard)

	for dto.StartDate.AddDate(0, 0, -1).Before(dto.EndDate) {
		date := dto.StartDate.In(s.Location).AddDate(0, 0, -1)
		dateKey := date.Format("2006-01-02")
		dateFormatted := date.Format("2006-01-02T15:04:05-0700")
		dateServicePointDashboardMap[dateKey] = &process_measures.ServicePointDashboard{
			Date:            dateFormatted,
			Periods:         periods,
			Magnitudes:      magnitudes,
			MagnitudeEnergy: magnitudeEnergy,
		}
		dto.StartDate = dto.StartDate.AddDate(0, 0, 1)
	}

	for _, dailyMonthly := range monthlyClosureMeasures {
		date := dailyMonthly.EndDate.In(s.Location).AddDate(0, 0, -1)
		dateKey := date.Format("2006-01-02")
		initDate := dailyMonthly.StartDate.In(s.Location).Format("2006-01-02T15:04:05-0700")
		monthlyMeasures := process_measures.MeasureServicePointDashboardMonthly{
			InitDate: initDate,
			ID:       dailyMonthly.Id,
			EndDate:  dailyMonthly.EndDate.In(s.Location).Format("2006-01-02T15:04:05-0700"),
			P0:       &measures.ValuesMonthly{},
			P1:       &measures.ValuesMonthly{},
			P2:       &measures.ValuesMonthly{},
			P3:       &measures.ValuesMonthly{},
			P4:       &measures.ValuesMonthly{},
			P5:       &measures.ValuesMonthly{},
			P6:       &measures.ValuesMonthly{},
		}

		for _, period := range append(allPeriods) {

			dailyMonthlyValuesPeriod := dailyMonthly.CalendarPeriods.GetPeriodValues(period)

			if dailyMonthlyValuesPeriod == nil {
				continue
			}

			measureMonthlyValues := measures.ValuesMonthly{
				Values: measures.Values{
					AI: dailyMonthlyValuesPeriod.AI,
					AE: dailyMonthlyValuesPeriod.AE,
					R1: dailyMonthlyValuesPeriod.R1,
					R2: dailyMonthlyValuesPeriod.R2,
					R3: dailyMonthlyValuesPeriod.R3,
					R4: dailyMonthlyValuesPeriod.R4,
				},
				AIi: dailyMonthlyValuesPeriod.AIi,
				AEi: dailyMonthlyValuesPeriod.AEi,
				R1i: dailyMonthlyValuesPeriod.R1i,
				R2i: dailyMonthlyValuesPeriod.R2i,
				R3i: dailyMonthlyValuesPeriod.R3i,
				R4i: dailyMonthlyValuesPeriod.R4i,
			}
			monthlyMeasures.SetPeriodMeasure(period, &measureMonthlyValues)
			monthlyMeasures.SetStatus(dailyMonthlyValuesPeriod.ValidationStatus)
		}

		servicePointDashboard := dateServicePointDashboardMap[dateKey]
		servicePointDashboard.SetMonthly(monthlyMeasures)

	}

	for _, dailyMeasure := range dailyClosureMeasures {
		date := dailyMeasure.EndDate.In(s.Location).AddDate(0, 0, -1)
		dateKey := date.Format("2006-01-02")
		dashboardMeasures := process_measures.MeasureServicePointDashboard{
			P0: &measures.Values{},
			P1: &measures.Values{},
			P2: &measures.Values{},
			P3: &measures.Values{},
			P4: &measures.Values{},
			P5: &measures.Values{},
			P6: &measures.Values{},
		}

		for _, period := range allPeriods {
			readingClosurePeriod := dailyMeasure.ToDailyReadingClosure().GetCalendarPeriod(period)
			if readingClosurePeriod == nil {
				continue
			}
			measure := measures.Values{
				AI: readingClosurePeriod.AI,
				AE: readingClosurePeriod.AE,
				R1: readingClosurePeriod.R1,
				R2: readingClosurePeriod.R2,
				R3: readingClosurePeriod.R3,
				R4: readingClosurePeriod.R4,
			}
			dashboardMeasures.SetPeriodMeasure(period, &measure)
			dashboardMeasures.SetStatus(readingClosurePeriod.ValidationStatus)
		}

		dateServicePoint := dateServicePointDashboardMap[dateKey]
		dateServicePoint.SetDaily(dashboardMeasures)

	}

	for _, curve := range curves {
		dateKey := curve.EndDate.In(s.Location).Add(-time.Hour * 1).Format("2006-01-02")
		if dashboardMeasure, ok := curveMap[dateKey]; ok {
			periodMeasure := dashboardMeasure.GetPeriodMeasure(curve.Period)
			for _, magnitude := range measures.ValidMagnitudes {
				magnitudeValue := curve.GetMagnitude(magnitude)
				periodMeasure.SumMagnitude(magnitude, magnitudeValue)
			}
			dashboardMeasure.SetPeriodMeasure(curve.Period, periodMeasure)
			if curve.Origin == measures.Filled {
				curve.ValidationStatus = measures.None
			}
			dashboardMeasure.SetStatusProcessCurve(curve.ValidationStatus)
			continue
		}

		servicePointMeasure := process_measures.MeasureServicePointDashboard{
			P0: &measures.Values{},
			P1: &measures.Values{},
			P2: &measures.Values{},
			P3: &measures.Values{},
			P4: &measures.Values{},
			P5: &measures.Values{},
			P6: &measures.Values{},
		}

		if curve.Origin == measures.Filled {
			servicePointMeasure.Status = measures.None
		} else {
			servicePointMeasure.Status = curve.ValidationStatus
		}

		servicePointMeasure.SetPeriodMeasure(curve.Period, &measures.Values{
			AI: curve.AI,
			AE: curve.AE,
			R1: curve.R1,
			R2: curve.R2,
			R3: curve.R3,
			R4: curve.R4,
		})

		curveMap[dateKey] = &servicePointMeasure
	}

	for date, dashboard := range curveMap {
		dateServicePoint := dateServicePointDashboardMap[date]
		dateServicePoint.SetCurves(*dashboard)
	}

	for _, value := range dateServicePointDashboardMap {
		result = append(result, *value)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	if err != nil {
		return process_measures.ServicePointDashboardWithType{}, err
	}

	span.SetAttributes(attribute.Int("response", len(result)))

	return process_measures.ServicePointDashboardWithType{
		Type:          meterConfig.Type,
		ServicePoints: result,
	}, nil
}
