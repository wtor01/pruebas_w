package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"math"
	"time"
)

type GetDashboardMeasures struct {
	inventoryRepository measures.InventoryRepository
	dashboardRepository process_measures.ProcessMeasureDashboardRepository
	tracer              trace.Tracer
	Location            *time.Location
}

func NewGetDashboardMeasures(dashboardRepository process_measures.ProcessMeasureDashboardRepository, inventoryRepository measures.InventoryRepository, loc *time.Location) *GetDashboardMeasures {
	return &GetDashboardMeasures{
		inventoryRepository: inventoryRepository,
		dashboardRepository: dashboardRepository,
		tracer:              telemetry.GetTracer(),
		Location:            loc,
	}
}

type GetDashboardMeasuresDTO struct {
	DistributorID string
	StartDate     time.Time
	EndDate       time.Time
}

func (g GetDashboardMeasuresDTO) String() string {
	return fmt.Sprintf("DistributorId: %s\nDates: %s - %s\n", g.DistributorID, g.StartDate.Format("2006-01-02 15-04"), g.EndDate.Format("2006-01-02 15-04"))
}

func NewGetDashboardMeasuresDTO(distributorID string, startDate time.Time, endDate time.Time) GetDashboardMeasuresDTO {
	endDate = endDate.AddDate(0, 0, 1)
	return GetDashboardMeasuresDTO{DistributorID: distributorID, StartDate: startDate, EndDate: endDate}
}

func (svc GetDashboardMeasures) getMeasureShouldBeInfo(ctx context.Context, dto GetDashboardMeasuresDTO) (measures.MeasureShouldBe, error) {
	equipments, err := svc.inventoryRepository.GroupMetersByType(ctx, measures.GroupMetersByTypeQuery{DistributorId: dto.DistributorID, StartDate: dto.StartDate, EndDate: dto.EndDate})

	if err != nil {
		return measures.MeasureShouldBe{}, err
	}

	diff := dto.EndDate.Sub(dto.StartDate)

	days := math.Ceil(diff.Hours() / 24)

	var result measures.MeasureShouldBe

	for t, curveTypes := range equipments {
		var dailyCurves int
		var dailyCloses int
		var monthlyCloses int
		for curveType, measureConfig := range curveTypes {
			dailyCloses += measureConfig.Count
			if curveType == measures.NoneType {
				continue
			}
			dailyCurves += measures.CalcByCurveType(curveType, measureConfig.Count)
			monthlyCloses += measureConfig.Total
		}

		totalCurves := int(days) * dailyCurves
		totalCloses := int(days) * dailyCloses

		if t == measures.TLG {
			result.Telegestion = struct {
				Curva   measures.MeasureShouldBeValue
				Closing measures.MeasureShouldBeValue
				Resumen measures.MeasureShouldBeValue
			}{
				Curva: measures.MeasureShouldBeValue{
					Daily: dailyCurves,
					Total: totalCurves,
				},
				Closing: measures.MeasureShouldBeValue{
					Daily: dailyCloses,
					Total: totalCloses,
				},
				Resumen: measures.MeasureShouldBeValue{
					Daily: 0,
					Total: monthlyCloses,
				},
			}
			continue
		}

		if t == measures.TLM {
			result.Telemedida = struct {
				Curva   measures.MeasureShouldBeValue
				Closing measures.MeasureShouldBeValue
			}{
				Curva: measures.MeasureShouldBeValue{
					Daily: dailyCurves,
					Total: totalCurves,
				},
				Closing: measures.MeasureShouldBeValue{
					Daily: dailyCloses,
					Total: totalCloses,
				},
			}
			continue
		}

		if t == measures.OTHER {
			result.Others = struct{ Closing measures.MeasureShouldBeValue }{
				Closing: measures.MeasureShouldBeValue{
					Daily: dailyCloses,
					Total: totalCloses,
				}}
			continue
		}
	}

	return result, nil
}

func (svc GetDashboardMeasures) Handle(ctx context.Context, dto GetDashboardMeasuresDTO) (measures.DashboardResult, error) {
	ctx, span := svc.tracer.Start(ctx, "GetDashboard_ProcessMeasures - Handle")
	defer span.End()
	span.SetAttributes(attribute.Stringer("params", dto))
	dashboardMeasures, err := svc.dashboardRepository.GetDashboard(ctx, process_measures.GetDashboardQuery{
		DistributorId: dto.DistributorID,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
	})

	now := time.Now().UTC()
	if !dto.EndDate.Before(now) {
		dto.EndDate = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, svc.Location).UTC()
	}

	if err != nil {
		return measures.DashboardResult{}, err
	}

	measureShouldBeInfo, err := svc.getMeasureShouldBeInfo(ctx, dto)

	if err != nil {
		return measures.DashboardResult{}, err
	}

	datesMeasures := measures.GetDashboardResultDailyData(dto.StartDate.In(svc.Location), dto.EndDate.In(svc.Location), measureShouldBeInfo)

	result := measures.GetDashboardResult(dashboardMeasures, measureShouldBeInfo, datesMeasures)

	return result, nil
}
