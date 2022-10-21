package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	scheduler "bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"context"
)

type UpdateScheduler struct {
	repository             process_measures.SchedulerRepository
	schedulerClientCreator scheduler.ClientCreator
	topic                  string
}

func NewUpdateScheduler(repository process_measures.SchedulerRepository, schedulerClientCreator scheduler.ClientCreator, topic string) *UpdateScheduler {
	return &UpdateScheduler{repository: repository, schedulerClientCreator: schedulerClientCreator, topic: topic}
}

type UpdateSchedulerDTO struct {
	ID            string
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
	Scheduler     string
	Description   string
}

func (s UpdateScheduler) Handler(ctx context.Context, dto UpdateSchedulerDTO) (process_measures.Scheduler, error) {

	sc, err := s.repository.GetScheduler(ctx, dto.ID)

	if err != nil {
		return process_measures.Scheduler{}, err
	}
	oldScheduler := sc.Clone()

	err = sc.Update(
		dto.Description,
		dto.DistributorId,
		dto.ServiceType,
		dto.PointType,
		dto.MeterType,
		dto.ReadingType,
		dto.Scheduler,
	)

	if err != nil {
		return process_measures.Scheduler{}, err
	}

	savedSchedulers, err := s.repository.SearchScheduler(ctx, process_measures.SearchScheduler{
		DistributorId: sc.DistributorId,
		ServiceType:   sc.ServiceType,
		PointType:     sc.PointType,
		MeterType:     sc.MeterType,
		ReadingType:   string(sc.ReadingType),
	})

	if err != nil {
		return process_measures.Scheduler{}, err
	}

	if (len(savedSchedulers) == 1 && savedSchedulers[0].ID != sc.ID) || len(savedSchedulers) > 1 {
		return process_measures.Scheduler{}, process_measures.ErrSchedulerExist
	}

	clientScheduler, err := s.schedulerClientCreator(ctx)

	defer clientScheduler.Close()

	err = clientScheduler.UpdateJob(ctx, &sc, s.topic)

	if err != nil {
		return process_measures.Scheduler{}, err
	}

	err = s.repository.SaveScheduler(ctx, sc)

	if err != nil {
		_ = clientScheduler.UpdateJob(ctx, &oldScheduler, s.topic)
		return process_measures.Scheduler{}, err
	}

	return sc, err
}
