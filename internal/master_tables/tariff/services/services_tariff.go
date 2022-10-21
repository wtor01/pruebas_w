package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
)

type TariffServices struct {
	GetTariffs           *GetTariffsService
	InsertTariff         *InsertTariffService
	DeleteTariff         *DeleteTariffService
	GetOneTariff         *GetOneTariffService
	ModifyTariff         *ModifyTariffsService
	GetAllTariffCalendar *GetTariffCalendarService
}

func NewTariffServices(repositoryTariff tariff.RepositoryTariff) *TariffServices {
	getGeographicZones := NewGetTariffRepositoryService(repositoryTariff)
	insertGeographicZones := NewInsertTariffService(repositoryTariff)
	deleteGeographicZones := NewDeleteTariffsService(repositoryTariff)
	getOneGeographicZone := NewGetOneTariffService(repositoryTariff)
	modifyGeographicZone := NewModifyTariffService(repositoryTariff)
	getAllTariffCalendar := NewGetTariffCalendarService(repositoryTariff)
	return &TariffServices{
		GetTariffs:           &getGeographicZones,
		InsertTariff:         &insertGeographicZones,
		DeleteTariff:         &deleteGeographicZones,
		GetOneTariff:         &getOneGeographicZone,
		ModifyTariff:         &modifyGeographicZone,
		GetAllTariffCalendar: &getAllTariffCalendar,
	}
}
