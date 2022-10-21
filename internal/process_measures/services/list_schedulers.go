package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
)

type ListSchedulerDto struct {
	Limit  int
	Offset *int
}

type ListScheduler struct {
	repository process_measures.SchedulerRepository
}

func NewListSchedulerService(repository process_measures.SchedulerRepository) *ListScheduler {
	return &ListScheduler{repository: repository}
}

func (s ListScheduler) Handler(ctx context.Context, dto ListSchedulerDto) ([]process_measures.Scheduler, int, error) {

	result, count, err := s.repository.ListScheduler(ctx, db.Pagination{
		Limit:  dto.Limit,
		Offset: dto.Offset,
	})

	if err != nil {
		return []process_measures.Scheduler{}, 0, err
	}

	return result, count, err
}
