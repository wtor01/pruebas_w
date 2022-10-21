package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
)

type SearchScheduler struct {
	repository process_measures.SchedulerRepository
}

func NewSearchSchedulerService(repository process_measures.SchedulerRepository) *SearchScheduler {
	return &SearchScheduler{repository: repository}
}

type SearchSchedulerDTO struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
}

func (s SearchScheduler) Handler(ctx context.Context, dto SearchSchedulerDTO) ([]process_measures.Scheduler, error) {

	searchSchedulers, err := s.repository.SearchScheduler(ctx, process_measures.SearchScheduler{
		DistributorId: dto.DistributorId,
		ServiceType:   dto.ServiceType,
		PointType:     dto.PointType,
		MeterType:     dto.MeterType,
		ReadingType:   dto.ReadingType,
	})

	if err != nil {
		return []process_measures.Scheduler{}, err
	}
	return searchSchedulers, nil

}
