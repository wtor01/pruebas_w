package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"time"
)

type ServicePointConfigView struct {
	MeterConfigId                 string    `gorm:"column:meter_config_id"`
	MeterId                       string    `gorm:"column:meter_id"`
	DistributorId                 string    `gorm:"column:distributor_id"`
	DistributorCode               string    `gorm:"column:distributor_code"`
	Cups                          string    `gorm:"column:cups"`
	MeterSerialNumber             string    `gorm:"column:meter_serial_number"`
	MetersType                    string    `gorm:"column:meters_type"`
	ServiceType                   string    `gorm:"column:service_type"`
	ReadingType                   string    `gorm:"column:reading_type"`
	PointType                     string    `gorm:"column:point_type"`
	MeasurePointType              string    `gorm:"column:measure_point_type"`
	PriorityContract              string    `gorm:"column:priority_contract"`
	TlgCode                       string    `gorm:"column:tlg_code"`
	TlgType                       string    `gorm:"column:tlg_type"`
	Calendar                      string    `gorm:"column:calendar"`
	StartDate                     time.Time `gorm:"column:start_date"`
	EndDate                       time.Time `gorm:"column:end_date"`
	CurveType                     string    `gorm:"column:curve_type"`
	Technology                    string    `gorm:"column:technology"`
	ContractualSituationId        string    `gorm:"column:contractual_situation_id"`
	TariffID                      string    `gorm:"column:tariff_id"`
	Coefficient                   string    `gorm:"column:coefficient"`
	CalendarCode                  string    `gorm:"column:calendar_code"`
	P1Demand                      float64   `gorm:"column:p1_demand"`
	P2Demand                      float64   `gorm:"column:p2_demand"`
	P3Demand                      float64   `gorm:"column:p3_demand"`
	P4Demand                      float64   `gorm:"column:p4_demand"`
	P5Demand                      float64   `gorm:"column:p5_demand"`
	P6Demand                      float64   `gorm:"column:p6_demand"`
	AI                            int       `gorm:"column:ai"`
	AE                            int       `gorm:"column:ae"`
	R1                            int       `gorm:"column:r1"`
	R2                            int       `gorm:"column:r2"`
	R3                            int       `gorm:"column:r3"`
	R4                            int       `gorm:"column:r4"`
	ContractualSituationStartDate time.Time `gorm:"column:contractual_situation_start_date"`
	ContractualSituationEndDate   time.Time `gorm:"column:contractual_situation_end_date"`
	ServicePointID                string    `gorm:"column:service_point_id"`
	GeographicID                  string    `gorm:"column:geographic_id"`
}

func (ServicePointConfigView) TableName() string {
	return "service_point_config_view"
}

func servicePointConfigViewToDomain(s ServicePointConfigView) inventory.ServicePointScheduler {
	return inventory.ServicePointScheduler{
		MeterConfigId:          s.MeterConfigId,
		MeterId:                s.MeterId,
		DistributorId:          s.DistributorId,
		DistributorCode:        s.DistributorCode,
		Cups:                   s.Cups,
		MeterSerialNumber:      s.MeterSerialNumber,
		ServiceType:            s.ServiceType,
		ReadingType:            s.ReadingType,
		PointType:              s.PointType,
		MeasurePointType:       s.MeasurePointType,
		PriorityContract:       s.PriorityContract,
		TlgCode:                s.TlgCode,
		TlgType:                s.TlgType,
		Calendar:               s.Calendar,
		Type:                   s.MetersType,
		StartDate:              s.StartDate,
		EndDate:                s.EndDate,
		CurveType:              s.CurveType,
		Technology:             s.Technology,
		ContractualSituationId: s.ContractualSituationId,
		TariffID:               s.TariffID,
		Coefficient:            s.Coefficient,
		CalendarCode:           s.CalendarCode,
		P1Demand:               s.P1Demand,
		P2Demand:               s.P2Demand,
		P3Demand:               s.P3Demand,
		P4Demand:               s.P4Demand,
		P5Demand:               s.P5Demand,
		P6Demand:               s.P6Demand,
		AI:                     s.AI,
		AE:                     s.AE,
		R1:                     s.R1,
		R2:                     s.R2,
		R3:                     s.R3,
		R4:                     s.R4,
		ContractStartDate:      s.ContractualSituationStartDate,
		ContractEndDate:        s.ContractualSituationEndDate,
		ServicePointID:         s.ServicePointID,
		GeographicID:           s.GeographicID,
	}
}

func (d InventoryPostgres) refreshServicePointConfigView() {
	d.db.Raw(`REFRESH MATERIALIZED VIEW service_point_config_view;`)
}

func (d InventoryPostgres) GetServicePoints(ctx context.Context, servicePointDto inventory.ServicePointSchedulerDto) ([]inventory.ServicePointScheduler, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var a []ServicePointConfigView
	dateFormat := servicePointDto.Time.UTC().Format("02/01/2006")

	tx := d.db.WithContext(ctxTimeout).Debug().Where("distributor_id = ?", servicePointDto.DistributorId).Where("meters_type IN (?)", servicePointDto.MeterType).Where("service_type = ?", servicePointDto.ServiceType).Where("point_type = ?", servicePointDto.PointType).Where("start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("end_date IS NULL or end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_end_date IS NULL or contractual_situation_end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat)

	if servicePointDto.Limit > 0 {
		tx = tx.Limit(servicePointDto.Limit).Offset(servicePointDto.Offset)
	}

	tx.Order("meter_id")

	tx.Find(&a)

	if tx.Error != nil {
		return []inventory.ServicePointScheduler{}, tx.Error
	}

	return utils.MapSlice(a, servicePointConfigViewToDomain), nil
}

func (d InventoryPostgres) GetServicePointByMeter(ctx context.Context, query inventory.GetMeasureEquipmentByMeterQuery) (inventory.ServicePointScheduler, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var result []ServicePointConfigView
	dateFormat := query.Time.Format("02/01/2006")

	q := d.db.WithContext(ctxTimeout).Debug().Where("start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("end_date IS NULL or end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_end_date IS NULL or contractual_situation_end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat)

	switch {
	case query.MeterSerialNumber != "":
		q.Where("meter_serial_number = ?", query.MeterSerialNumber)

	case query.MeterID != "":
		q.Where("meter_id = ?", query.MeterID)

	default:
		return inventory.ServicePointScheduler{}, nil
	}

	tx := q.Limit(1).Find(&result)

	if tx.Error != nil {
		return inventory.ServicePointScheduler{}, tx.Error
	}

	if len(result) == 0 {
		return inventory.ServicePointScheduler{}, nil
	}

	return servicePointConfigViewToDomain(result[0]), nil
}

func (d InventoryPostgres) GetServicePointByCups(ctx context.Context, query inventory.GetServicePointByCupsQuery) (inventory.ServicePointScheduler, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var result []ServicePointConfigView
	dateFormat := query.Time.Format("02/01/2006")

	q := d.db.WithContext(ctxTimeout).Debug().Where("start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("end_date IS NULL or end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_end_date IS NULL or contractual_situation_end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat)

	if query.CUPS != "" {
		q.Where("cups = ?", query.CUPS)
	} else {
		return inventory.ServicePointScheduler{}, nil

	}
	tx := q.Limit(1).Find(&result)

	if tx.Error != nil {
		return inventory.ServicePointScheduler{}, tx.Error
	}

	if len(result) == 0 {
		return inventory.ServicePointScheduler{}, nil
	}

	return servicePointConfigViewToDomain(result[0]), nil
}
