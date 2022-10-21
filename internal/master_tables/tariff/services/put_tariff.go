package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"context"
)

type ModifyTariffsService struct {
	TariffRepository tariff.RepositoryTariff
}

//NewModifyTariffService create new Modify service
func NewModifyTariffService(geographicRepository tariff.RepositoryTariff) ModifyTariffsService {
	return ModifyTariffsService{TariffRepository: geographicRepository}
}

//Handler handle modify service
func (s ModifyTariffsService) Handler(ctx context.Context, code string, tf tariff.Tariffs) error {
	err := s.TariffRepository.ModifyTariff(ctx, code, tf)
	if err != nil {
		return err
	}
	err = s.TariffRepository.InsertTariffCalendar(ctx, tf.CalendarId, tf, tf.UpdatedBy)
	return err
}
