package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"context"
)

type CreateScheduler struct {
	repository             billing_measures.BillingSchedulerRepository
	schedulerClientCreator scheduler.ClientCreator
	topic                  string
}

func NewCreateScheduler(repository billing_measures.BillingSchedulerRepository, schedulerClientCreator scheduler.ClientCreator, topic string) *CreateScheduler {
	return &CreateScheduler{repository: repository, schedulerClientCreator: schedulerClientCreator, topic: topic}
}

type CreateSchedulerDTO struct {
	ID            string
	Name          string
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ProcessType   string
	Scheduler     string
}

func (s CreateScheduler) Handler(ctx context.Context, dto CreateSchedulerDTO) (billing_measures.Scheduler, error) {
	sc, err := billing_measures.NewScheduler(
		dto.ID,
		dto.Name,
		"",
		dto.DistributorId,
		dto.ServiceType,
		dto.PointType,
		dto.MeterType,
		dto.ProcessType,
		dto.Scheduler)

	if err != nil {
		return billing_measures.Scheduler{}, err
	}
	savedSchedulers, err := s.repository.SearchScheduler(ctx, billing_measures.SearchScheduler{
		DistributorId: sc.DistributorId,
		ServiceType:   sc.ServiceType,
		PointType:     sc.PointType,
		MeterType:     sc.MeterType,
		ProcessType:   sc.ProcessType,
	})

	if err != nil {
		return billing_measures.Scheduler{}, err
	}

	if len(savedSchedulers) != 0 {
		return billing_measures.Scheduler{}, billing_measures.ErrSchedulerExist
	}

	clientScheduler, err := s.schedulerClientCreator(ctx)

	defer clientScheduler.Close()

	jobId, err := clientScheduler.CreateJob(ctx, &sc, s.topic)
	if err != nil {
		return billing_measures.Scheduler{}, err
	}
	sc.SetSchedulerId(jobId)

	err = s.repository.SaveScheduler(ctx, sc)

	if err != nil {
		clientScheduler.DeleteJob(ctx, sc.SchedulerId)
		return billing_measures.Scheduler{}, err
	}

	return sc, err

}
