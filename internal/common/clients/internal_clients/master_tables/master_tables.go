package master_tables

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
)

type MasterTables struct {
	TariffServices *services.TariffServices
}

func NewMasterTables(tariffServices *services.TariffServices) *MasterTables {
	return &MasterTables{TariffServices: tariffServices}
}

func (m MasterTables) GetTariff(ctx context.Context, dto clients.GetTariffDto) (clients.Tariffs, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetTariff client")
	defer span.End()

	tariff, err := m.TariffServices.GetOneTariff.Handler(ctx, services.GetTariffByCodeDto{
		ID: dto.ID,
	})

	return clients.Tariffs{
		Id:           tariff.Id,
		Description:  tariff.Description,
		TensionLevel: tariff.TensionLevel,
		CodeOdos:     tariff.CodeOdos,
		CodeOne:      tariff.CodeOne,
		Periods:      tariff.Periods,
		GeographicId: tariff.GeographicId,
		CalendarId:   tariff.CalendarId,
		Coef:         tariff.Coef,
	}, err
}
