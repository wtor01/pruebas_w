package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"context"
)

type InsertTariffService struct {
	TariffRepository tariff.RepositoryTariff
}

//NewInsertTariffService create a service struct
func NewInsertTariffService(geographicRepository tariff.RepositoryTariff) InsertTariffService {
	return InsertTariffService{TariffRepository: geographicRepository}
}

//Handler handle insert service
func (s InsertTariffService) Handler(ctx context.Context, tf tariff.Tariffs) error {
	//convert to uppercase
	err := s.TariffRepository.InsertTariff(ctx, tf)
	if err != nil {
		return err
	}
	err = s.TariffRepository.InsertTariffCalendar(ctx, tf.CalendarId, tf, tf.CreatedBy)
	return err
}
