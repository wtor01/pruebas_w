package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
)

type InsertDayTypeService struct {
	dayTypesRepository seasons.RepositorySeasons
}

func NewInsertDayTypeService(dayTypesRepository seasons.RepositorySeasons) InsertDayTypeService {
	return InsertDayTypeService{dayTypesRepository: dayTypesRepository}
}

func (s InsertDayTypeService) Handler(ctx context.Context, id string, dayT seasons.DayTypes) error {
	err := s.dayTypesRepository.InsertDayTypes(ctx, id, dayT)

	return err
}
