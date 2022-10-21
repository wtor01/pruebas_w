package billing_measures

import (
	"context"
	"time"
)

type PointSelfConsumptionType string

const (
	FronteraServicePointType PointSelfConsumptionType = "Frontera"
	GdServicePointType       PointSelfConsumptionType = "G-D"
	ConsumoServicePointType  PointSelfConsumptionType = "Consumo"
	SsaaServicePointType     PointSelfConsumptionType = "SSAA"
)

type ConnectionType string

const (
	InternalConnection ConnectionType = "Red Interior"
	ExternalConnection ConnectionType = "Próxima a través de red"
)

type PointSelfConsumption struct {
	ID               string                   `json:"_id"`
	CauID            string                   `json:"cau_ID"`
	ServicePointType PointSelfConsumptionType `json:"service_point_type"`
	CUPS             string                   `json:"CUPS"`
	InitDate         time.Time                `json:"init_date"`
	EndDate          time.Time                `json:"end_date"`
	InstalationFlag  int                      `json:"instalation_flag"`
	WithoutmeterFlag int                      `json:"withoutmeter_flag"`
	Exent1Flag       int                      `json:"exent1_flag"`
	Exent2Flag       int                      `json:"exent2_flag"`
	PartitionCoeff   float64                  `json:"partition_coeff"`
	DistributorId    string                   `json:"distributor_id"`
}

type ConfigSelfConsumption struct {
	ID                    string         `json:"_id"`
	CauID                 string         `json:"cau_ID"`
	StatusID              int            `json:"status_ID"`
	StatusName            string         `json:"status_name"`
	InitDate              time.Time      `json:"init_date"`
	EndDate               time.Time      `json:"end_date"`
	CreationDate          time.Time      `json:"creation_date"`
	DistributorId         string         `json:"distributor_id"`
	CnmcTypeId            int            `json:"cnmc_type_id"`
	CnmcTypeName          string         `json:"cnmc_type_name"`
	CnmcTypeDesc          string         `json:"cnmc_type_desc"`
	ConfType              ConfigType     `json:"conf_type"`
	ConfTypeDescription   string         `json:"conf_type_description"`
	ConsumerType          string         `json:"consumer_type"`
	ParticipantNumber     int            `json:"participant_number"`
	ConnType              ConnectionType `json:"conn_type"`
	Excedents             bool           `json:"excedents"`
	Compensation          bool           `json:"compensation"`
	GenerationPot         float64        `json:"generation_pot"`
	GroupSubgroup         int            `json:"group_subgroup"`
	AntivertType          string         `json:"antivert_type"`
	SolarZoneId           int            `json:"solar_zone_id"`
	SolarZoneNum          int            `json:"solar_zone_num"`
	SolarZoneName         string         `json:"solar_zone_name"`
	TechnologyId          string         `json:"technology_id"`
	TechnologyDescription string         `json:"technology_description"`
}

type SelfConsumption struct {
	ID            string                  `json:"_id"`
	CAU           string                  `json:"CAU"`
	Name          string                  `json:"name"`
	StatusID      int                     `json:"status_ID"`
	StatusName    string                  `json:"status_Name"`
	CcaaId        int                     `json:"ccaaId"`
	Ccaa          string                  `json:"ccaa"`
	InitDate      time.Time               `json:"init_date"`
	EndDate       time.Time               `json:"end_date"`
	DistributorId string                  `json:"distributor_id"`
	Configs       []ConfigSelfConsumption `json:"configs"`
	Points        []PointSelfConsumption  `json:"points"`
}

func (s SelfConsumption) GetBillingSelfConsumptionConfig(billingMeasure BillingMeasure) BillingSelfConsumptionConfig {
	for _, c := range s.Configs {
		if c.InitDate.After(billingMeasure.InitDate) || c.EndDate.Before(billingMeasure.EndDate) {
			continue
		}
		return BillingSelfConsumptionConfig{
			Id:                    c.ID,
			CauID:                 c.CauID,
			StatusID:              c.StatusID,
			StatusName:            c.StatusName,
			InitDate:              c.InitDate,
			EndDate:               c.EndDate,
			CnmcTypeId:            c.CnmcTypeId,
			CnmcTypeName:          c.CnmcTypeName,
			CnmcTypeDesc:          c.CnmcTypeDesc,
			ConfType:              c.ConfType,
			ConfTypeDescription:   c.ConfTypeDescription,
			ConsumerType:          c.ConsumerType,
			ParticipantNumber:     c.ParticipantNumber,
			ConnType:              c.ConnType,
			Excedents:             c.Excedents,
			Compensation:          c.Compensation,
			GenerationPot:         c.GenerationPot,
			GroupSubgroup:         c.GroupSubgroup,
			AntivertType:          c.AntivertType,
			SolarZoneId:           c.SolarZoneId,
			SolarZoneNum:          c.SolarZoneNum,
			SolarZoneName:         c.SolarZoneName,
			TechnologyId:          c.TechnologyId,
			TechnologyDescription: c.TechnologyDescription,
		}
	}

	return BillingSelfConsumptionConfig{}
}

/*
	BillingSelftConsumption API
*/
func NewBillingSelfConsumptionUnit(consumption BillingSelfConsumption) BillingSelfConsumptionUnit {
	return BillingSelfConsumptionUnit{
		EndDate:  consumption.EndDate,
		InitDate: consumption.InitDate,
		Status:   consumption.Status,
		CauInfo: BillingSelfConsumptionCau{
			Id:         consumption.CAU,
			Name:       consumption.Name,
			Points:     len(consumption.Points),
			UnitType:   consumption.Config.CnmcTypeDesc,
			ConfigType: consumption.Config.ConfType,
		},
	}
}

type BillingSelfConsumptionUnit struct {
	EndDate             time.Time                                   `json:"end_date"`
	InitDate            time.Time                                   `json:"init_date"`
	Status              BillingSelfConsumptionStatus                `json:"status"`
	CauInfo             BillingSelfConsumptionCau                   `json:"cau_info" json:"cau_info"`
	Totals              BillingSelfConsumptionTotals                `json:"totals" json:"totals"`
	NetGeneration       []BillingSelfConsumptionNetGeneration       `json:"net_generation" json:"net_generation,omitempty"`
	UnitConsumption     []BillingSelfConsumptionUnitConsumption     `json:"unit_consumption" json:"unit_consumption,omitempty"`
	CalendarConsumption []BillingSelfConsumptionCalendarConsumption `json:"calendar_consumption" json:"calendar_consumption,omitempty"`
	Cups                []BillingSelfConsumptionCups                `json:"cups" json:"cups,omitempty"`
}

func (b *BillingSelfConsumptionUnit) SumTotals(values BillingSelfConsumptionTotals) {
	b.Totals.sum(values)
}

func (b *BillingSelfConsumptionUnit) SetNetGeneration(netGeneration []BillingSelfConsumptionNetGeneration) {
	b.NetGeneration = netGeneration
}

func (b *BillingSelfConsumptionUnit) SetUnitConsumption(unitConsumption []BillingSelfConsumptionUnitConsumption) {
	b.UnitConsumption = unitConsumption
}

func (b *BillingSelfConsumptionUnit) SetCalendarConsumption(calendarConsumption []BillingSelfConsumptionCalendarConsumption) {
	b.CalendarConsumption = calendarConsumption
}

func (b *BillingSelfConsumptionUnit) SetCups(cupsList []BillingSelfConsumptionCups) {
	b.Cups = cupsList
}

type BillingSelfConsumptionCau struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Points     int        `json:"points"`
	UnitType   string     `json:"unit_type"`
	ConfigType ConfigType `json:"config_type"`
}

type BillingSelfConsumptionTotals struct {
	GrossGeneration    float64 `json:"gross_generation"`
	NetGeneration      float64 `json:"net_generation"`
	SelfConsumption    float64 `json:"self_consumption"`
	NetworkConsumption float64 `json:"network_consumption"`
	AuxConsumption     float64 `json:"aux_consumption"`
}

func (totals *BillingSelfConsumptionTotals) sum(values BillingSelfConsumptionTotals) {
	totals.GrossGeneration += values.GrossGeneration
	totals.NetGeneration += values.NetGeneration
	totals.SelfConsumption += values.SelfConsumption
	totals.NetworkConsumption += values.NetworkConsumption
	totals.AuxConsumption += values.AuxConsumption
}

type BillingSelfConsumptionNetGeneration struct {
	Date      string  `json:"date"`
	Net       float64 `json:"net"`
	Excedente float64 `json:"excedente"`
}

func (n *BillingSelfConsumptionNetGeneration) Sum(net, excedente float64) {
	n.Net += net
	n.Excedente += excedente
}

type BillingSelfConsumptionUnitConsumption struct {
	Date    string  `json:"date"`
	Network float64 `json:"network"`
	Self    float64 `json:"self"`
	Aux     float64 `json:"aux"`
}

func (n *BillingSelfConsumptionUnitConsumption) Sum(network, self, aux float64) {
	n.Network += network
	n.Self += self
	n.Aux += aux
}

type BillingSelfConsumptionCalendarConsumption struct {
	Date   string                      `json:"date"`
	Energy float64                     `json:"energy"`
	Values []CalendarConsumptionValues `json:"values"`
}

func (c *BillingSelfConsumptionCalendarConsumption) AddValue(calendarValue CalendarConsumptionValues) {
	c.Values = append(c.Values, calendarValue)
}

func (c *BillingSelfConsumptionCalendarConsumption) Sum(energy float64) {
	c.Energy += energy
}

type BillingSelfConsumptionCups struct {
	Cups        string                   `json:"cups"`
	Type        PointSelfConsumptionType `json:"type"`
	Consumption float64                  `json:"consumption"`
	Generation  float64                  `json:"generation"`
	StartDate   string                   `json:"start_date"`
	EndDate     string                   `json:"end_date"`
}

func (c *BillingSelfConsumptionCups) Sum(consumption, generation float64) {
	c.Generation += generation
	c.Consumption += consumption
}

type CalendarConsumptionValues struct {
	Hour   string  `json:"hour"`
	Energy float64 `json:"energy"`
}

type QueryGetActiveSelfConsumptionByCUP struct {
	CUP       string
	StartDate time.Time
	EndDate   time.Time
}

type GetSelfConsumptionByCUP struct {
	DistributorId string
	CUP           string
	Date          time.Time
}

type GetSelfConsumptionByDistributortDto struct {
	DistributorId string
	Limit         int
	Offset        int
	Date          time.Time
	Values        map[string]string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=SelfConsumptionRepository
type SelfConsumptionRepository interface {
	GetActiveSelfConsumptionByCUP(ctx context.Context, query QueryGetActiveSelfConsumptionByCUP) (SelfConsumption, error)
	GetSelfConsumptionByCUP(ctx context.Context, query GetSelfConsumptionByCUP) (SelfConsumption, error)
	GetSelfConsumptionActiveByDistributor(ctx context.Context, query GetSelfConsumptionByDistributortDto) ([]SelfConsumption, int, error)
}
