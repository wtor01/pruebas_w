package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
)

func (c MeterConfig) toDomain() inventory.MeterConfig {
	return inventory.MeterConfig{
		ID:               c.ID.String(),
		StartDate:        c.StartDate,
		EndDate:          c.EndDate,
		AI:               c.AI,
		AE:               c.AE,
		R1:               c.R1,
		R2:               c.R2,
		R3:               c.R3,
		R4:               c.R4,
		M:                c.M,
		E:                c.E,
		CurveType:        c.CurveType,
		ReadingType:      c.ReadingType,
		CalendarID:       c.CalendarID,
		PriorityContract: c.PriorityContract,
		MeterID:          c.MeterID,
		MeasurePointID:   c.MeasurePointID,
		ServicePointID:   c.ServicePointID,
		TlgCode:          measures.MeterConfigTlgCode(c.TlgCode),
		Type:             c.Type,
		RentingPrice:     c.RentingPrice,
		SerialNumber:     c.SerialNumber,
		DistributorID:    c.DistributorID.String(),
		CUPS:             c.CUPS,
		MeasurePointType: measures.MeasurePointType(c.MeasurePointType),
	}
}

func (d InventoryPostgres) GetMeterConfig(ctx context.Context, search inventory.GetMeterConfigByMeterQuery) (inventory.MeterConfig, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	config := MeterConfig{}

	tx := d.db.WithContext(ctxTimeout).First(&config, "serial_number = ? AND start_date > ? AND end_date < ?", search.MeterSerialNumber, search.Date, search.Date)

	return config.toDomain(), tx.Error
}
