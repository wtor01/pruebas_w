package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
)

//GetSeasonByIdDto struct send code for Get
type GetSeasonByIdDto struct {
	Id string
}
type GetSeasonByIdService struct {
	SeasonRepository seasons.RepositorySeasons
}

//NewGetOneCalendarService create get one struct
func NewGetSeasonByIdService(seasonRepository seasons.RepositorySeasons) GetSeasonByIdService {
	return GetSeasonByIdService{SeasonRepository: seasonRepository}
}

//Handler handle get one geographic service
func (s GetSeasonByIdService) Handler(ctx context.Context, dto GetSeasonByIdDto) (seasons.Seasons, error) {
	res, err := s.SeasonRepository.GetSeason(ctx, dto.Id)
	return res, err
}
