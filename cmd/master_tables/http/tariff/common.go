package tariff

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
)

func TariffToResponse(d tariff.Tariffs) Tariffs {
	return Tariffs{
		Id:           d.Id,
		TariffCode:   d.Id,
		CodeOdos:     d.CodeOdos,
		CodeOne:      d.CodeOne,
		Description:  d.Description,
		GeographicId: d.GeographicId,
		Periods:      d.Periods,
		TensionLevel: TariffsTensionLevel(d.TensionLevel),
		CalendarCode: d.CalendarId,
		Coef:         TariffsCoef(d.Coef),
	}
}
