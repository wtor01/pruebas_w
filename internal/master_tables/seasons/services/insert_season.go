package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
)

type InsertSeasonService struct {
	SeasonRepository seasons.RepositorySeasons
}

func NewInsertSeasonService(seasonRepository seasons.RepositorySeasons) InsertSeasonService {
	return InsertSeasonService{SeasonRepository: seasonRepository}
}

//Handler handle insert service
func (s InsertSeasonService) Handler(ctx context.Context, season seasons.Seasons) error {
	//convert to uppercase
	err := s.SeasonRepository.InsertSeason(ctx, season)
	return err
}
