package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"context"
)

type GetAllDto struct {
	Q           string
	Limit       int
	Offset      *int
	Sort        map[string]string
	CurrentUser auth.User
}

type GetGeographicsService struct {
	geographicRepository geographic.RepositoryGeographic
}

// NewGetGeographicRepositoryService Generate Get distributor service
func NewGetGeographicRepositoryService(geographicRepository geographic.RepositoryGeographic) GetGeographicsService {
	return GetGeographicsService{geographicRepository: geographicRepository}
}

// Handler Get data, count and errors from db
func (s GetGeographicsService) Handler(ctx context.Context, dto GetAllDto) ([]geographic.GeographicZones, int, error) {
	ds, count, err := s.geographicRepository.GetAllGeographicZones(ctx, geographic.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
