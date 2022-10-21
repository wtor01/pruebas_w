package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
	uuid "github.com/satori/go.uuid"
)

type DeleteSeasonByIdDto struct {
	Id uuid.UUID
}

type DeleteSeasonService struct {
	SeasonRepository seasons.RepositorySeasons
}

func NewDeleteSeasonService(seasonRepository seasons.RepositorySeasons) DeleteSeasonService {
	return DeleteSeasonService{SeasonRepository: seasonRepository}
}

func (s DeleteSeasonService) Handler(ctx context.Context, dto DeleteSeasonByIdDto) error {
	err := s.SeasonRepository.DeleteSeason(ctx, dto.Id)
	return err
}
