package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"context"
	"strings"
)

type InsertGeographicsService struct {
	GeographicRepository geographic.RepositoryGeographic
}

//NewInsertGeographicsService create a service struct
func NewInsertGeographicsService(geographicRepository geographic.RepositoryGeographic) InsertGeographicsService {
	return InsertGeographicsService{GeographicRepository: geographicRepository}
}

//Handler handle insert service
func (s InsertGeographicsService) Handler(ctx context.Context, zones geographic.GeographicZones) error {
	//convert to uppercase
	zones.Code = strings.ToUpper(zones.Code)
	err := s.GeographicRepository.InsertGeographicZone(ctx, zones)

	return err
}
