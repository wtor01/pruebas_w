package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
)

type GetDashboardMeasuresStats struct {
	repository process_measures.ProcessMeasureDashboardStatsRepository
	tracer     trace.Tracer
}

func NewGetDashboardMeasuresStats(repository process_measures.ProcessMeasureDashboardStatsRepository) *GetDashboardMeasuresStats {
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

func (s GetDashboardMeasuresStats) Handler(ctx context.Context, dto GetDashboardMeasuresStatsDTO) ([]process_measures.ProcessMeasureDashboardStatsGlobal, error) {
	ctx, span := s.tracer.Start(ctx, "GetDashboardMeasuresStats - Handler")
	defer span.End()

	if dto.Year < 0 || dto.Month < 0 || dto.Type == "" || dto.DistributorID == "" {
		return []process_measures.ProcessMeasureDashboardStatsGlobal{}, errors.New("invalid input format")
	}
	stats, err := s.repository.GetStatisticsGlobal(ctx, process_measures.SearchDashboardStats{
		DistributorID: dto.DistributorID,
		Month:         dto.Month,
		Year:          dto.Year,
		Type:          dto.Type,
	})
	if err != nil {
		return []process_measures.ProcessMeasureDashboardStatsGlobal{}, err
	}
	return stats, err
}
