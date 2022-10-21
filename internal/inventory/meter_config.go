package inventory

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

type MeterConfig struct {
	ID               string
	StartDate        time.Time
	EndDate          time.Time
	AI               int
	AE               int
	R1               int
	R2               int
	R3               int
	R4               int
	M                int
	E                int
	CurveType        string
	ReadingType      string
	CalendarID       string
	PriorityContract string
	MeterID          string
	MeasurePointID   string
	ServicePointID   string
	TlgCode          measures.MeterConfigTlgCode
	Type             string
	RentingPrice     float64
	SerialNumber     string
	DistributorID    string
	CUPS             string
	MeasurePointType measures.MeasurePointType
}
