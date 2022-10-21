package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"time"
)

type ListClosureDto struct {
	Id            string
	DistributorId string
	Cups          string
	Moment        string
	StartDate     time.Time
	EndDate       time.Time
}

type GetClosure struct {
	Location        *time.Location
	mongoRepository process_measures.ProcessMeasureClosureRepository
}

func NewGetClosureService(mongoRepository process_measures.ProcessMeasureClosureRepository) *GetClosure {
	return &GetClosure{mongoRepository: mongoRepository}
}

func (s GetClosure) Handler(ctx context.Context, dto ListClosureDto) (process_measures.ProcessedMonthlyClosure, error) {

	var processClosure process_measures.GetClosure

	processClosure = process_measures.GetClosure{
		DistributorId: dto.DistributorId,
	}

	if dto.Id != "" {

		processClosure.Id = dto.Id
		processClosure.Moment = process_measures.SelectMoment(dto.Moment)

	} else {
		processClosure.CUPS = dto.Cups
		processClosure.StartDate = dto.StartDate
		processClosure.EndDate = dto.EndDate
	}

	result, err := s.mongoRepository.GetClosure(ctx, processClosure)

	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}

	return result, nil

}
