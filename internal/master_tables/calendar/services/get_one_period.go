package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"context"
	uuid "github.com/satori/go.uuid"
)

//GetPeriodByIdDto struct send code for Get
type GetPeriodByIdDto struct {
	Id uuid.UUID
}
type GetOnePeriodService struct {
	PeriodRepository calendar.RepositoryCalendar
}

//NewGetOnePeriodService get one struct
func NewGetOnePeriodService(periodRepository calendar.RepositoryCalendar) GetOnePeriodService {
	return GetOnePeriodService{PeriodRepository: periodRepository}
}

//Handler handle get one geographic service
func (s GetOnePeriodService) Handler(ctx context.Context, dto GetPeriodByIdDto) (calendar.PeriodCalendar, error) {
	res, err := s.PeriodRepository.GetPeriodById(ctx, dto.Id)
	return res, err
}
