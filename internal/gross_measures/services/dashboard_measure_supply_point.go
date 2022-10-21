package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/async"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.opentelemetry.io/otel/trace"
	"sort"
	"time"
)

type DashboardMeasureSupplyPointDto struct {
	Cups          string
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

func NewDashboardMeasureSupplyPointDto(cups, distributorId string, startDate, endDate time.Time) DashboardMeasureSupplyPointDto {
	return DashboardMeasureSupplyPointDto{
		Cups:          cups,
		DistributorId: distributorId,
		StartDate:     startDate,
		EndDate:       endDate.AddDate(0, 0, 1),
	}
}

type DashboardMeasureSupplyPointService struct {
	grossRepository     gross_measures.GrossMeasureRepository
	inventoryRepository measures.InventoryRepository
	calendarRepository  measures.CalendarPeriodRepository
	masterTablesClient  clients.MasterTables
	location            *time.Location
	tracer              trace.Tracer
}

func NewDashboardMeasureSupplyPointService(grossRepository gross_measures.GrossMeasureRepository, inventoryRepository measures.InventoryRepository, calendarRepository measures.CalendarPeriodRepository, masterTableClients clients.MasterTables, loc *time.Location) *DashboardMeasureSupplyPointService {
	return &DashboardMeasureSupplyPointService{
		grossRepository:     grossRepository,
		inventoryRepository: inventoryRepository,
		calendarRepository:  calendarRepository,
		masterTablesClient:  masterTableClients,
		location:            loc,
		tracer:              telemetry.GetTracer(),
	}
}

func (s DashboardMeasureSupplyPointService) Handler(ctx context.Context, dto DashboardMeasureSupplyPointDto) (gross_measures.DashboardMeasureSupplyPoint, error) {
	ctx, span := s.tracer.Start(ctx, "DashboardMeasureSupplyPointService - Handler")
	defer span.End()

	meterConfig, periods, err := s.getMeterInfo(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        dto.Cups,
		Time:        dto.EndDate,
		Distributor: dto.DistributorId,
	})

	if err != nil {
		return gross_measures.DashboardMeasureSupplyPoint{}, err
	}

	grossCurves, grossCloses, err := s.getGrossMeasures(ctx, gross_measures.QueryListMeasure{
		SerialNumber: meterConfig.SerialNumber(),
		StartDate:    dto.StartDate,
		EndDate:      dto.EndDate,
	})

	if err != nil {
		return gross_measures.DashboardMeasureSupplyPoint{}, err
	}

	dashboardMeasures := gross_measures.NewDashboardMeasureSupplyPoint(meterConfig, periods)

	periodsMap := map[measures.PeriodKey]struct{}{
		measures.P0: {},
	}

	for _, period := range dashboardMeasures.Periods {
		periodsMap[period] = struct{}{}
	}

	calendarDailyClosureMap := make(map[string]gross_measures.CalendarStatus)
	calendarMonthlyClosureMap := make(map[string]gross_measures.CalendarStatus)
	calendarCurveMap := make(map[string]gross_measures.CalendarStatus)
	listDailyCloses := make(map[string]gross_measures.ListDailyClose)
	listMonthlyCloses := make(map[string]gross_measures.ListMonthlyClose)

	for _, curve := range grossCurves {
		subTime := time.Hour
		if curve.CurveType == measures.QuarterMeasureCurveReadingType {
			subTime = time.Minute * 15
		}
		dateFormatted := s.formatDate(curve.EndDate, subTime)
		s.addToCalendar(calendarCurveMap, gross_measures.CalendarStatus{
			Date:   dateFormatted,
			Status: curve.Status,
		})
	}
	dashboardMeasures.SetCalendar(measures.Curve, s.transformCalendars(calendarCurveMap))

	for _, grossClose := range grossCloses {
		dateFormatted := s.formatDate(grossClose.EndDate, time.Hour)
		calendarStatus := gross_measures.CalendarStatus{
			Date:   dateFormatted,
			Status: grossClose.Status,
		}

		if grossClose.ReadingType == measures.BillingClosure {
			s.addToCalendar(calendarMonthlyClosureMap, calendarStatus)
			s.addToMonthlyCloseList(listMonthlyCloses, grossClose, calendarStatus, periodsMap)
			continue
		}
		s.addToCalendar(calendarDailyClosureMap, calendarStatus)
		s.addToDailyCloseList(listDailyCloses, grossClose, calendarStatus, periodsMap)

	}
	dashboardMeasures.SetCalendar(measures.DailyClosure, s.transformCalendars(calendarDailyClosureMap))
	dashboardMeasures.SetCalendar(measures.BillingClosure, s.transformCalendars(calendarMonthlyClosureMap))
	dashboardMeasures.SetListDailyClosure(s.transformDailyCloses(listDailyCloses))
	dashboardMeasures.SetListMonthlyClosure(s.transformMonthlyCloses(listMonthlyCloses))

	return dashboardMeasures, nil
}

/*
	Get info of meter (Config, tariff & calendar periods)
*/
func (s DashboardMeasureSupplyPointService) getMeterInfo(ctx context.Context, query measures.GetMeterConfigByCupsQuery) (measures.MeterConfig, []measures.PeriodKey, error) {
	meter, err := s.inventoryRepository.GetMeterConfigByCups(ctx, query)

	if err != nil {
		return measures.MeterConfig{}, []measures.PeriodKey{}, err
	}

	tariff, err := s.masterTablesClient.GetTariff(ctx, clients.GetTariffDto{
		ID: meter.TariffID(),
	})

	if err != nil {
		return measures.MeterConfig{}, []measures.PeriodKey{}, err
	}

	calendar, err := s.calendarRepository.GetCalendarPeriod(ctx, measures.SearchCalendarPeriod{
		Day:          query.Time,
		GeographicID: tariff.GeographicId,
		CalendarCode: tariff.CalendarId,
		Location:     s.location,
	})

	return meter, calendar.GetAllPeriods(), err
}

/*
	Get gross measures (Curve & Closes) in date range
*/
func (s DashboardMeasureSupplyPointService) getGrossMeasures(ctx context.Context, query gross_measures.QueryListMeasure) ([]gross_measures.MeasureCurveWrite, []gross_measures.MeasureCloseWrite, error) {

	asyncCloses := async.Exec(ctx, func(ctx context.Context) ([]gross_measures.MeasureCloseWrite, error) {
		return s.grossRepository.ListCloseMeasures(ctx, query)
	})
	asyncCurves := async.Exec(ctx, func(ctx context.Context) ([]gross_measures.MeasureCurveWrite, error) {
		return s.grossRepository.ListCurveMeasures(ctx, query)
	})

	curves := asyncCurves.Await(ctx)
	closes := asyncCloses.Await(ctx)

	if curves.Error != nil {
		return []gross_measures.MeasureCurveWrite{}, []gross_measures.MeasureCloseWrite{}, curves.Error
	}
	if closes.Error != nil {
		return []gross_measures.MeasureCurveWrite{}, []gross_measures.MeasureCloseWrite{}, closes.Error
	}

	return curves.Result, closes.Result, nil
}

func (s DashboardMeasureSupplyPointService) formatDate(date time.Time, subTime time.Duration) string {
	return date.In(s.location).Add(-subTime).Format("2006-01-02")
}

/*
	Functions to add into maps
*/
func (s DashboardMeasureSupplyPointService) addToCalendar(items map[string]gross_measures.CalendarStatus, calendarStatus gross_measures.CalendarStatus) {
	calendarValue, ok := items[calendarStatus.Date]

	if !ok {
		items[calendarStatus.Date] = calendarStatus
		return
	}

	calendarValue.SetStatus(calendarStatus.Status)
	items[calendarStatus.Date] = calendarValue
}

func (s DashboardMeasureSupplyPointService) addToDailyCloseList(items map[string]gross_measures.ListDailyClose, grossDaily gross_measures.MeasureCloseWrite, calendarStatus gross_measures.CalendarStatus, periodsMap map[measures.PeriodKey]struct{}) {
	mapKey := grossDaily.GetOriginFile()
	dailyClose := gross_measures.ListDailyClose{
		EndDate: calendarStatus.Date,
		Status:  calendarStatus.Status,
		Origin:  mapKey,
	}

	for _, periodValue := range grossDaily.Periods {
		if _, ok := periodsMap[periodValue.Period]; !ok {
			continue
		}

		measureValues := measures.Values{
			AI: periodValue.AI,
			AE: periodValue.AE,
			R1: periodValue.R1,
			R2: periodValue.R2,
			R3: periodValue.R3,
			R4: periodValue.R4,
		}

		dailyClose.SetPeriodValues(periodValue.Period, measureValues)
	}
	items[mapKey] = dailyClose
}

func (s DashboardMeasureSupplyPointService) addToMonthlyCloseList(items map[string]gross_measures.ListMonthlyClose, grossMonthly gross_measures.MeasureCloseWrite, calendarStatus gross_measures.CalendarStatus, periodsMap map[measures.PeriodKey]struct{}) {
	mapKey := grossMonthly.GetOriginFile()
	_, ok := items[mapKey]

	monthlyClose := gross_measures.ListMonthlyClose{
		EndDate:  calendarStatus.Date,
		Status:   calendarStatus.Status,
		InitDate: s.formatDate(grossMonthly.StartDate, time.Hour),
		Origin:   mapKey,
	}

	if ok {
		monthlyClose = items[mapKey]
	}

	for _, periodValue := range grossMonthly.Periods {
		if _, okPeriod := periodsMap[periodValue.Period]; !okPeriod {
			continue
		}
		measureValues := measures.Values{
			AI: periodValue.AI,
			AE: periodValue.AE,
			R1: periodValue.R1,
			R2: periodValue.R2,
			R3: periodValue.R3,
			R4: periodValue.R4,
		}

		if grossMonthly.Type == measures.Incremental {
			monthlyClose.SetPeriodMonthlyValues(periodValue.Period, measureValues)
			continue
		}
		monthlyClose.SetPeriodValues(periodValue.Period, measureValues)
	}

	items[mapKey] = monthlyClose
}

/*
	Functions to transform maps in slices
*/
func (s DashboardMeasureSupplyPointService) transformMonthlyCloses(items map[string]gross_measures.ListMonthlyClose) []gross_measures.ListMonthlyClose {

	monthlyCloses := utils.MapToSlice(items)

	sort.Slice(monthlyCloses, func(i, j int) bool {
		monthlyClose, compareClose := monthlyCloses[i], monthlyCloses[j]

		if monthlyClose.EndDate != compareClose.EndDate {
			return monthlyCloses[i].EndDate < monthlyCloses[j].EndDate

		}

		return !monthlyClose.Status.Compare(measures.StatusValue, compareClose.Status)
	})

	return monthlyCloses
}

func (s DashboardMeasureSupplyPointService) transformDailyCloses(items map[string]gross_measures.ListDailyClose) []gross_measures.ListDailyClose {

	dailyCloses := utils.MapToSlice(items)

	sort.Slice(dailyCloses, func(i, j int) bool {
		dailyClose, compareClose := dailyCloses[i], dailyCloses[j]

		if dailyClose.EndDate != compareClose.EndDate {
			return dailyCloses[i].EndDate < dailyCloses[j].EndDate

		}

		return !dailyClose.Status.Compare(measures.StatusValue, compareClose.Status)
	})

	return dailyCloses
}

func (s DashboardMeasureSupplyPointService) transformCalendars(items map[string]gross_measures.CalendarStatus) []gross_measures.CalendarStatus {

	calendars := utils.MapToSlice(items)

	sort.Slice(calendars, func(i, j int) bool {
		return calendars[i].Date < calendars[j].Date
	})

	return calendars
}
