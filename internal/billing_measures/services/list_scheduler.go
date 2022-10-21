package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
)

type ListSchedulerDto struct {
	Limit  int
	Offset *int
}
type ListScheduler struct {
	repository billing_measures.BillingSchedulerRepository
}

func NewListSchedulerService(repository billing_measures.BillingSchedulerRepository) *ListScheduler {
	return &ListScheduler{repository: repository}
}
func (s ListScheduler) Handler(ctx context.Context, dto ListSchedulerDto) ([]billing_measures.Scheduler, int, error) {

	result, count, err := s.repository.ListScheduler(ctx, db.Pagination{
		Limit:  dto.Limit,
		Offset: dto.Offset,
	})

	if err != nil {
		return []billing_measures.Scheduler{}, 0, err
	}

	return result, count, err
}
