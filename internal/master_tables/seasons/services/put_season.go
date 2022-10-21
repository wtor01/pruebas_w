package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
	uuid "github.com/satori/go.uuid"
)

type ModifySeasonService struct {
	SeasonRepository seasons.RepositorySeasons
}

func NewModifySeasonService(seasonRepository seasons.RepositorySeasons) ModifySeasonService {
	return ModifySeasonService{SeasonRepository: seasonRepository}
}

func (s ModifySeasonService) Handler(ctx context.Context, id uuid.UUID, season seasons.Seasons) error {
	err := s.SeasonRepository.ModifySeason(ctx, id, season)
	return err
}
