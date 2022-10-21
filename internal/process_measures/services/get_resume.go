package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"time"
)

type ListResumeDto struct {
	DistributorId string
	Cups          string
	StartDate     time.Time
	EndDate       time.Time
}

type GetResume struct {
	Location        *time.Location
	mongoRepository process_measures.ProcessMeasureClosureRepository
}

func NewGetResumeService(mongoRepository process_measures.ProcessMeasureClosureRepository, location *time.Location) *GetResume {
	return &GetResume{mongoRepository: mongoRepository, Location: location}
}

// Handler se trae el cierre anterior y el cierre posterior y lo devuelve
func (s GetResume) Handler(ctx context.Context, dto ListResumeDto) (process_measures.ResumesProcessMonthlyClosure, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetResume")
	defer span.End()
	processResume := process_measures.GetResume{
		DistributorId: dto.DistributorId,
		Cups:          dto.Cups,
		StartDate:     dto.StartDate.In(s.Location),
		EndDate:       dto.EndDate.In(s.Location).AddDate(0, 0, 1),
	}

	result, err := s.mongoRepository.GetResume(ctx, processResume)

	if err != nil {
		return process_measures.ResumesProcessMonthlyClosure{}, err
	}

	return result, nil

}
