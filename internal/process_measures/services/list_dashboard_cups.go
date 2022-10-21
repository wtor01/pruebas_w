package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type ListDashboardCupsService struct {
	inventoryRepository measures.InventoryRepository
	Location            *time.Location
	tracer              trace.Tracer
	dashboardRepository process_measures.ProcessMeasureDashboardRepository
}

func NewListDashboardCupsDTO(distributorId string, limit int64, offset int64, meterType string, startDate time.Time, endDate time.Time) ListDashboardCupsDTO {
	endDate = endDate.AddDate(0, 0, 1)
	return ListDashboardCupsDTO{DistributorId: distributorId, Limit: limit, Offset: offset, Type: meterType, StartDate: startDate, EndDate: endDate}
}

type ListDashboardCupsDTO struct {
	DistributorId string
	Type          string
	Limit         int64
	Offset        int64
	StartDate     time.Time
	EndDate       time.Time
}

func (l ListDashboardCupsDTO) String() string {
	return fmt.Sprintf("DistributorId: %s\nType: %s\nDates: %s - %s\nLimit: %d\nOffset: %d", l.DistributorId, l.Type, l.StartDate.Format("2006-01-02 15-04"), l.EndDate.Format("2006-01-02 15-04"), l.Limit, l.Offset)
}

func NewListDashboardCupsService(inventoryRepository measures.InventoryRepository, dashboardRepository process_measures.ProcessMeasureDashboardRepository, location *time.Location) *ListDashboardCupsService {
	return &ListDashboardCupsService{
		inventoryRepository: inventoryRepository,
		dashboardRepository: dashboardRepository,
		tracer:              telemetry.GetTracer(),
		Location:            location,
	}
}

func (s ListDashboardCupsService) Handle(ctx context.Context, query ListDashboardCupsDTO) (measures.DashboardListCups, error) {
	ctx, span := s.tracer.Start(ctx, "ListDashboardCups_ProcessMeasures - Handle")
	defer span.End()
	result, err := s.inventoryRepository.GetMetersAndCountByDistributorId(ctx, measures.GetMetersAndCountByDistributorIdQuery{
		DistributorId: query.DistributorId,
		Type:          measures.MeterType(query.Type),
		StartDate:     query.StartDate,
		EndDate:       query.EndDate,
		Limit:         query.Limit,
		Offset:        query.Offset,
	})

	if err != nil {
		return measures.DashboardListCups{}, err
	}

	cupsList, err := s.dashboardRepository.GetCupsMeasures(ctx, process_measures.ListCupsQuery{
		Cups: utils.MapSlice(result.Data, func(item measures.GetMetersAndCountData) string {
			return item.Cups
		}),
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	})

	if err != nil {
		return measures.DashboardListCups{}, err
	}

	dashboardCups := s.transformToDashboardCups(cupsList, result.Data, measures.MeterType(query.Type), query.StartDate, query.EndDate)

	span.SetAttributes(attribute.Stringer("params", query))
	span.SetAttributes(attribute.Int("cupsList", len(dashboardCups)))
	return measures.DashboardListCups{
		Cups:  dashboardCups,
		Total: result.Count,
	}, nil
}

func (s ListDashboardCupsService) transformToDashboardCups(cupsList map[string]*measures.DashboardCupsReading, metersData []measures.GetMetersAndCountData, meterType measures.MeterType, startDate, endDate time.Time) []measures.DashboardCups {
	dashboardCups := make([]measures.DashboardCups, 0, cap(metersData))

	for _, data := range metersData {

		s.setShouldBe(cupsList[data.Cups], data.Meters, meterType, startDate, endDate)
		dashboardCups = append(dashboardCups, measures.DashboardCups{
			Cups:   data.Cups,
			Values: *cupsList[data.Cups],
		})
	}
	return dashboardCups
}

func (s ListDashboardCupsService) setShouldBe(cupsReading *measures.DashboardCupsReading, meters []measures.GetMetersAndCountMeter, meterType measures.MeterType, startDate, endDate time.Time) {

	now := time.Now().UTC()

	if now.Before(endDate) {
		endDate = now
	}

	diffDays := utils.DiffDays(startDate, endDate)

	for _, readingType := range measures.ReadingTypes {
		var shouldBe int
		switch readingType {
		case measures.Curve:
			if meterType == measures.OTHER {
				break
			}
			restDays := diffDays
			initDate := startDate
			for _, meter := range meters {
				if !meter.EndDate.Before(endDate) {
					shouldBe += measures.CalcByCurveType(meter.CurveType, 1) * restDays
					break
				}
				diff := utils.DiffDays(initDate, meter.EndDate)
				restDays -= diff
				initDate = meter.EndDate
				shouldBe += measures.CalcByCurveType(meter.CurveType, 1) * diff
			}
		case measures.DailyClosure:
			if meterType == measures.OTHER {
				break
			}
			shouldBe = diffDays
		case measures.BillingClosure:
			for _, meter := range meters {
				monthClose := time.Date(meter.EndDate.Year(), meter.EndDate.Month(), 32, meter.EndDate.Hour(), meter.EndDate.Minute(), meter.EndDate.Second(), 0, time.UTC)
				if meter.EndDate.Before(endDate) && !monthClose.Equal(meter.EndDate) {
					shouldBe += 1
				}
			}
			for initDate := startDate.AddDate(0, 1, 0); initDate.Before(endDate); initDate = initDate.AddDate(0, 1, 0) {
				shouldBe += 1
			}
		}

		cupsReading.SetShouldBe(readingType, shouldBe)
	}

}
