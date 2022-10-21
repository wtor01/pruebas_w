package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/trace"
)

type GetClosureResumeDashboardDTO struct {
	BillingMeasureID string
}

type GetClosureResumeDashboard struct {
	repository billing_measures.BillingMeasuresDashboardRepository
	tracer     trace.Tracer
}

func NewGetClosureResumeDashboard(repository billing_measures.BillingMeasuresDashboardRepository) *GetClosureResumeDashboard {
	return &GetClosureResumeDashboard{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

func (closureResume GetClosureResumeDashboard) Handler(ctx context.Context, dto GetClosureResumeDashboardDTO) (billing_measures.BillingMeasureDashboardResumeClosure, error) {
	ctx, span := closureResume.tracer.Start(ctx, "SearchBillingMeasureClosureResume - Handler")
	defer span.End()

	bM, err := closureResume.repository.SearchBillingMeasureClosureResume(ctx, dto.BillingMeasureID)
	if err != nil {
		return billing_measures.BillingMeasureDashboardResumeClosure{}, err
	}
	toReturn := billing_measures.NewBillingMeasureDashboardResumeClosure(bM)
	return toReturn, err
}
