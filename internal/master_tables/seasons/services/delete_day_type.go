package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
	uuid "github.com/satori/go.uuid"
)

type DeleteDayTypeByIdDto struct {
	Id uuid.UUID
}

type DeleteDayTypeService struct {
	dayTypesRepository seasons.RepositorySeasons
}

func NewDeleteDayTypeService(dayTypesRepository seasons.RepositorySeasons) DeleteDayTypeService {
	return DeleteDayTypeService{dayTypesRepository: dayTypesRepository}
}

func (s DeleteDayTypeService) Handler(ctx context.Context, dto DeleteSeasonByIdDto) error {
	err := s.dayTypesRepository.DeleteDayType(ctx, dto.Id)
	return err
}
