package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
)

type SearchScheduler struct {
	repository billing_measures.BillingSchedulerRepository
}

func NewSearchSchedulerService(repository billing_measures.BillingSchedulerRepository) *SearchScheduler {
	return &SearchScheduler{repository: repository}
}

type SearchSchedulerDTO struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
}

func (s SearchScheduler) Handler(ctx context.Context, dto SearchSchedulerDTO) ([]billing_measures.Scheduler, error) {

	searchSchedulers, err := s.repository.SearchScheduler(ctx, billing_measures.SearchScheduler{
		DistributorId: dto.DistributorId,
		ServiceType:   dto.ServiceType,
		PointType:     dto.PointType,
		MeterType:     dto.MeterType,
	})

	if err != nil {
		return []billing_measures.Scheduler{}, err
	}
	return searchSchedulers, nil

}
