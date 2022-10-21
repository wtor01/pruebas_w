package inventory

import (
	"time"
)

type ServicePointScheduler struct {
	MeterConfigId          string
	MeterId                string
	DistributorId          string
	DistributorCode        string
	Cups                   string
	MeterSerialNumber      string
	ServiceType            string
	ReadingType            string
	PointType              string
	MeasurePointType       string
	PriorityContract       string
	TlgCode                string
	TlgType                string
	Calendar               string
	Type                   string
	StartDate              time.Time
	EndDate                time.Time
	CurveType              string
	Technology             string
	ContractualSituationId string
	TariffID               string
	Coefficient            string
	CalendarCode           string
	P1Demand               float64
	P2Demand               float64
	P3Demand               float64
	P4Demand               float64
	P5Demand               float64
	P6Demand               float64
	AI                     int
	AE                     int
	R1                     int
	R2                     int
	R3                     int
	R4                     int
	ContractStartDate      time.Time
	ContractEndDate        time.Time
	ServicePointID         string
	GeographicID           string
}

type ServicePointSchedulerDto struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
	Time          time.Time
	Limit         int
	Offset        int
}
