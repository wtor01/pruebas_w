package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.opentelemetry.io/otel/trace"
	"sort"
	"time"
)

type DashboardSupplyPointCurvesDto struct {
	Date          time.Time                        `json:"date"`
	DistributorId string                           `json:"distributor_id"`
	Cups          string                           `json:"cups"`
	CurveType     measures.MeasureCurveReadingType `json:"curve_type"`
}

func NewDashboardSupplyPointCurvesDto(cups, distributorId, curveType string, date time.Time) DashboardSupplyPointCurvesDto {
	return DashboardSupplyPointCurvesDto{
		Date:          date,
		DistributorId: distributorId,
		Cups:          cups,
		CurveType:     measures.MeasureCurveReadingType(curveType),
	}
}

type DashboardSupplyPointCurvesService struct {
	grossMeasureRepository gross_measures.GrossMeasureRepository
	inventoryRepository    measures.InventoryRepository
	tracer                 trace.Tracer
	location               *time.Location
}

func NewDashboardSupplyPointCurvesService(
	grossMeasureRepository gross_measures.GrossMeasureRepository,
	inventoryRepository measures.InventoryRepository,
	loc *time.Location,
) *DashboardSupplyPointCurvesService {
	return &DashboardSupplyPointCurvesService{
		grossMeasureRepository: grossMeasureRepository,
		inventoryRepository:    inventoryRepository,
		tracer:                 telemetry.GetTracer(),
		location:               loc,
	}
}

func (s *DashboardSupplyPointCurvesService) Handler(ctx context.Context, dto DashboardSupplyPointCurvesDto) ([]gross_measures.DashboardSupplyPointCurve, error) {
	ctx, span := s.tracer.Start(ctx, "DashboardSupplyPointCurvesService - Handler")
	defer span.End()

	meterConfig, err := s.inventoryRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        dto.Cups,
		Time:        dto.Date,
		Distributor: dto.DistributorId,
	})

	if err != nil {
		return nil, err
	}

	var curveType measures.MeasureCurveReadingType
	if meterConfig.Type == measures.TLM {
		curveType = dto.CurveType
	}

	curves, err := s.grossMeasureRepository.ListDailyCurveMeasures(ctx, gross_measures.QueryListForProcessCurve{
		SerialNumber: meterConfig.SerialNumber(),
		Date:         dto.Date,
		CurveType:    curveType,
	})

	if err != nil {
		return nil, err
	}

	curveMap := make(map[string]gross_measures.DashboardSupplyPointCurve)

	for _, curve := range curves {
		subTime := time.Hour
		if curve.CurveType == measures.QuarterMeasureCurveReadingType {
			subTime = time.Minute * 15
		}
		dateFormatted := curve.EndDate.In(s.location).Add(-subTime).Format("15:04")

		s.addToCurveMap(curveMap, gross_measures.DashboardSupplyPointCurve{
			Hour:   dateFormatted,
			Status: curve.Status,
			File:   curve.GetOriginFile(),
			Values: measures.Values{
				AI: curve.AI,
				AE: curve.AE,
				R1: curve.R1,
				R2: curve.R2,
				R3: curve.R3,
				R4: curve.R4,
			},
		})
	}

	return s.transformCurveMap(curveMap), nil
}

func (s DashboardSupplyPointCurvesService) addToCurveMap(items map[string]gross_measures.DashboardSupplyPointCurve, point gross_measures.DashboardSupplyPointCurve) {
	items[point.Hour] = point
}

func (s DashboardSupplyPointCurvesService) transformCurveMap(items map[string]gross_measures.DashboardSupplyPointCurve) []gross_measures.DashboardSupplyPointCurve {
	curves := utils.MapToSlice(items)

	sort.Slice(curves, func(i, j int) bool {
		return curves[i].Hour < curves[j].Hour
	})

	return curves
}
