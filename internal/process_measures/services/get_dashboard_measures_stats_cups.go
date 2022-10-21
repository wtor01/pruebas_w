package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
)

type GetDashboardMeasuresStatsCups struct {
	repository process_measures.ProcessMeasureDashboardStatsRepository
	tracer     trace.Tracer
}

func NewGetDashboardMeasuresStatsCups(repository process_measures.ProcessMeasureDashboardStatsRepository) *GetDashboardMeasuresStatsCups {
	return &GetDashboardMeasuresStatsCups{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

type GetDashboardMeasuresStatsCupsDTO struct {
	DistributorID string
	Month         int
	Year          int
	Type          string
	Offset        int
	Limit         int
}

func (s GetDashboardMeasuresStatsCups) Handler(ctx context.Context, dto GetDashboardMeasuresStatsCupsDTO) ([]process_measures.ProcessMeasureDashboardStatsGlobal, int, error) {
	ctx, span := s.tracer.Start(ctx, "GetDashboardMeasuresStatsCups - Handler")
	defer span.End()

	if dto.Year < 0 || dto.Month < 0 || dto.Type == "" || dto.DistributorID == "" {
		return []process_measures.ProcessMeasureDashboardStatsGlobal{}, 0, errors.New("invalid input format")
	}
	stats, cont, err := s.repository.GetStatisticsCups(ctx, process_measures.SearchDashboardStats{
		DistributorID: dto.DistributorID,
		Month:         dto.Month,
		Year:          dto.Year,
		Type:          dto.Type,
		Offset:        dto.Offset,
		Limit:         dto.Limit,
	})
	if err != nil {
		return []process_measures.ProcessMeasureDashboardStatsGlobal{}, 0, err
	}
	return stats, cont, err
}
