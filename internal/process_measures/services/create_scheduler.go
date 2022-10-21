package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	scheduler "bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"context"
)

type CreateScheduler struct {
	repository             process_measures.SchedulerRepository
	schedulerClientCreator scheduler.ClientCreator
	topic                  string
}

func NewCreateScheduler(repository process_measures.SchedulerRepository, schedulerClientCreator scheduler.ClientCreator, topic string) *CreateScheduler {
	return &CreateScheduler{repository: repository, schedulerClientCreator: schedulerClientCreator, topic: topic}
}

type CreateSchedulerDTO struct {
	ID            string
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
	Scheduler     string
	Name          string
	Description   string
}

func (s CreateScheduler) Handler(ctx context.Context, dto CreateSchedulerDTO) (process_measures.Scheduler, error) {

	sc, err := process_measures.NewScheduler(
		dto.ID,
		dto.Name,
		dto.Description,
		dto.DistributorId,
		"",
		dto.ServiceType,
		dto.PointType,
		dto.MeterType,
		dto.ReadingType,
		dto.Scheduler)

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

	if len(savedSchedulers) != 0 {
		return process_measures.Scheduler{}, process_measures.ErrSchedulerExist
	}

	clientScheduler, err := s.schedulerClientCreator(ctx)

	defer clientScheduler.Close()

	jobId, err := clientScheduler.CreateJob(ctx, &sc, s.topic)

	if err != nil {
		return process_measures.Scheduler{}, err
	}

	sc.SetSchedulerId(jobId)

	err = s.repository.SaveScheduler(ctx, sc)

	if err != nil {
		clientScheduler.DeleteJob(ctx, sc.SchedulerId)
		return process_measures.Scheduler{}, err
	}

	return sc, err
}
