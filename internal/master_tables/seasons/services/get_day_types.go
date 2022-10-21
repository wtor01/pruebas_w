package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
)

type GetDayTypesService struct {
	dayTypesRepository seasons.RepositorySeasons
}

func NewGetDayTypesService(dayTypesRepository seasons.RepositorySeasons) GetDayTypesService {
	return GetDayTypesService{dayTypesRepository: dayTypesRepository}
}

// Handler Get data, count and errors from db
func (s GetDayTypesService) Handler(ctx context.Context, seasonId string, dto GetAllDto) ([]seasons.DayTypes, int, error) {
	ds, count, err := s.dayTypesRepository.GetAllDayTypes(ctx, seasonId, seasons.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
