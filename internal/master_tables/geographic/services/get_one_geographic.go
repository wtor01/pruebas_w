package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"context"
)

//GetGeographicByIdDto struct send code for Get
type GetGeographicByIdDto struct {
	Id string
}
type GetOneGeographicsService struct {
	GeographicRepository geographic.RepositoryGeographic
}

//NewGetOneGeographicsService create get one struct
func NewGetOneGeographicsService(geographicRepository geographic.RepositoryGeographic) GetOneGeographicsService {
	return GetOneGeographicsService{GeographicRepository: geographicRepository}
}

//Handler handle get one geographic service
func (s GetOneGeographicsService) Handler(ctx context.Context, dto GetGeographicByIdDto) (geographic.GeographicZones, error) {
	res, err := s.GeographicRepository.GetGeographicZone(ctx, dto.Id)

	return res, err
}
