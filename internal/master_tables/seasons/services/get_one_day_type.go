package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
)

//GetDayTypeByIdDto struct send code for Get
type GetDayTypeByIdDto struct {
	Id string
}
type GetDayTypeByIdService struct {
	DayTypeRepository seasons.RepositorySeasons
}

//NewGetOneCalendarService create get one struct
func NewGetDayTypeByIdService(seasonRepository seasons.RepositorySeasons) GetDayTypeByIdService {
	return GetDayTypeByIdService{DayTypeRepository: seasonRepository}
}

//Handler handle get one geographic service
func (s GetDayTypeByIdService) Handler(ctx context.Context, dto GetDayTypeByIdDto) (seasons.DayTypes, error) {
	res, err := s.DayTypeRepository.GetDayType(ctx, dto.Id)
	return res, err
}
