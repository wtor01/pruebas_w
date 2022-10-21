package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
)

type GetDashboardMeasuresStats struct {
	repository gross_measures.GrossMeasuresDashboardStatsRepository
	tracer     trace.Tracer
}

func NewGetDashboardMeasuresStats(repository gross_measures.GrossMeasuresDashboardStatsRepository) *GetDashboardMeasuresStats {
	return &GetDashboardMeasuresStats{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

type GetDashboardMeasuresStatsDTO struct {
	DistributorID string
	Month         int
	Year          int
	Type          string
}

func (s GetDashboardMeasuresStats) Handler(ctx context.Context, dto GetDashboardMeasuresStatsDTO) ([]gross_measures.GrossMeasuresDashboardStatsGlobal, error) {
	ctx, span := s.tracer.Start(ctx, "GetGrossMeasuresDashboardStats - Handler")
	defer span.End()

	if dto.Year < 0 || (dto.Month < 0 && dto.Month > 12) || dto.Type == "" || dto.DistributorID == "" {
		return []gross_measures.GrossMeasuresDashboardStatsGlobal{}, errors.New("invalid input format")
	}
	stats, err := s.repository.GetStatisticsGlobal(ctx, gross_measures.SearchDashboardStats{
		DistributorID: dto.DistributorID,
		Month:         dto.Month,
		Year:          dto.Year,
		Type:          dto.Type,
	})
	if err != nil {
		return []gross_measures.GrossMeasuresDashboardStatsGlobal{}, err
	}
	return stats, err
}
