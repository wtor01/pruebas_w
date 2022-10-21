package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
)

type ListDashboardMeasuresStatsSerialNumber struct {
	repository gross_measures.GrossMeasuresDashboardStatsRepository
	tracer     trace.Tracer
}

func NewListDashboardMeasuresStatsSerialNumber(repository gross_measures.GrossMeasuresDashboardStatsRepository) *ListDashboardMeasuresStatsSerialNumber {
	return &ListDashboardMeasuresStatsSerialNumber{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

type ListDashboardMeasuresStatsSerialNumberDTO struct {
	DistributorId string
	Month         int
	Year          int
	Type          string
	Ghost         bool
	Offset        int
	Limit         int
}

func (s ListDashboardMeasuresStatsSerialNumber) Handler(ctx context.Context, dto ListDashboardMeasuresStatsSerialNumberDTO) (gross_measures.ListGrossMeasuresStatisticsSerialNumberResult, error) {
	ctx, span := s.tracer.Start(ctx, "ListDashboardMeasuresStatsSerialNumber - Handler")
	defer span.End()

	if dto.Month < 1 || dto.Month > 12 || dto.Type == "" || dto.DistributorId == "" || dto.Year < 0 {
		return gross_measures.ListGrossMeasuresStatisticsSerialNumberResult{
			Data:  []gross_measures.DashboardSerialNumber{},
			Count: 0,
		}, errors.New("not valid input")
	}

	result, err := s.repository.ListGrossMeasuresStatisticsSerialNumber(ctx, gross_measures.SearchDashboardSerialNumber{
		DistributorId: dto.DistributorId,
		Month:         dto.Month,
		Year:          dto.Year,
		Type:          dto.Type,
		Ghost:         dto.Ghost,
		Offset:        dto.Offset,
		Limit:         dto.Limit,
	})
	return result, err
}
