package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
)

//GeographicServices struct all services
type GeographicServices struct {
	GetGeographicZones    *GetGeographicsService
	InsertGeographicZone  *InsertGeographicsService
	DeleteGeographicsZone *DeleteGeographicsService
	GetOneGeographicsZone *GetOneGeographicsService
	ModifyGeographicZone  *ModifyGeographicsService
}

//NewGeographicServices create GeographicServices
func NewGeographicServices(repositoryGeographic geographic.RepositoryGeographic, publisher event.PublisherCreator, topic string) *GeographicServices {
	getGeographicZones := NewGetGeographicRepositoryService(repositoryGeographic)
	insertGeographicZones := NewInsertGeographicsService(repositoryGeographic)
	deleteGeographicZones := NewDeleteGeographicsService(repositoryGeographic, publisher, topic)
	getOneGeographicZone := NewGetOneGeographicsService(repositoryGeographic)
	modifyGeographicZone := NewModifyGeographicsService(repositoryGeographic, publisher, topic)
	return &GeographicServices{
		GetGeographicZones:    &getGeographicZones,
		InsertGeographicZone:  &insertGeographicZones,
		DeleteGeographicsZone: &deleteGeographicZones,
		GetOneGeographicsZone: &getOneGeographicZone,
		ModifyGeographicZone:  &modifyGeographicZone,
	}
}
