package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type DashboardSummaryService struct {
	dashboardRepository billing_measures.BillingMeasuresDashboardRepository
	tracer              trace.Tracer
}

type DashboardSummaryDto struct {
	DistributorId string
	MeterType     measures.MeterType
	StartDate     time.Time
	EndDate       time.Time
}

func NewDashboardSummaryDto(distributorId, meterType string, startDate, endDate time.Time) DashboardSummaryDto {
	return DashboardSummaryDto{
		DistributorId: distributorId,
		MeterType:     measures.MeterType(meterType),
		StartDate:     startDate,
		EndDate:       endDate.AddDate(0, 0, 1),
	}
}

func NewDashboardSummaryService(dashboardRepository billing_measures.BillingMeasuresDashboardRepository) *DashboardSummaryService {
	return &DashboardSummaryService{
		dashboardRepository: dashboardRepository,
		tracer:              telemetry.GetTracer(),
	}
}

func (s DashboardSummaryService) Handler(ctx context.Context, dto DashboardSummaryDto) (billing_measures.FiscalMeasureSummary, error) {
	ctx, span := s.tracer.Start(ctx, "DashboardSummaryService - Handler")
	defer span.End()

	span.SetAttributes(attribute.String("distributorId", dto.DistributorId))
	span.SetAttributes(attribute.String("meterType", string(dto.MeterType)))
	fiscalSummary, err := s.dashboardRepository.GroupFiscalMeasureSummary(ctx, billing_measures.GroupFiscalMeasureSummaryQuery{
		DistributorId: dto.DistributorId,
		MeterType:     dto.MeterType,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
	})

	return fiscalSummary, err
}
