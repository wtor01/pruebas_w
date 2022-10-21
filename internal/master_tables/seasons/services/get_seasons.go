package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"context"
)

type GetAllDto struct {
	Q           string
	Limit       int
	Offset      *int
	Sort        map[string]string
	CurrentUser auth.User
}

type GetSeasonsService struct {
	seasonsRepository seasons.RepositorySeasons
}

// NewGetSeasonsRepositoryService Generate Get distributor service
func NewGetSeasonsRepositoryService(seasonsRepository seasons.RepositorySeasons) GetSeasonsService {
	return GetSeasonsService{seasonsRepository: seasonsRepository}
}

// Handler Get data, count and errors from db
func (s GetSeasonsService) Handler(ctx context.Context, dto GetAllDto) ([]seasons.Seasons, int, error) {
	ds, count, err := s.seasonsRepository.GetAllSeasons(ctx, seasons.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
