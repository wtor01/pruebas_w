package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
	uuid "github.com/satori/go.uuid"
)

type ModifyDayTypeService struct {
	dayTypeRepository seasons.RepositorySeasons
}

func NewModifyDayTypeService(dayTypeRepository seasons.RepositorySeasons) ModifyDayTypeService {
	return ModifyDayTypeService{dayTypeRepository: dayTypeRepository}
}

func (s ModifyDayTypeService) Handler(ctx context.Context, seasonId uuid.UUID, sea seasons.DayTypes) error {
	err := s.dayTypeRepository.ModifyDayTypes(ctx, seasonId, sea)
	return err
}
