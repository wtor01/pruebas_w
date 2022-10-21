package calendar

import (
	"github.com/satori/go.uuid"
	"time"
)

type Calendar struct {
	Id             string
	Code           string
	Description    string
	Periods        int
	GeographicCode string
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
	UpdatedAt      time.Time
}

type DayType string

const (
	Festive  DayType = "Festive"
	Workable DayType = "Workable"
)

//PeriodCalendar service struct
type PeriodCalendar struct {
	ID             uuid.UUID
	CalendarCode   string
	GeographicCode string
	PeriodNumber   string
	Description    string
	Year           int
	DayType        DayType
	StartHour      int
	EndHour        int
	StartDate      string
	EndDate        string
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
	Energy         bool
	Power          bool
	UpdatedAt      time.Time
}

func (period PeriodCalendar) IsFestive() bool {
	return period.DayType == Festive
}
