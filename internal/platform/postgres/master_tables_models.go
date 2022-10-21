package postgres

import (
	"time"
)

type GeographicZones struct {
	ModelEntity
	Code        string `gorm:"code" gorm:"index"`
	Description string `gorm:"description"`
	ModelRegisterUserActions
}

type Calendars struct {
	ID           string          `gorm:"id" gorm:"primaryKey" gorm:"index"`
	Description  string          `gorm:"description"`
	Periods      int             `gorm:"periods"`
	GeographicID *string         `gorm:"column:geographic_id;type:uuid"  gorm:"foreignKey:GeographicId;references:ID"`
	Geographic   GeographicZones `gorm:"geographic_id"`
	CreatedAt    time.Time       `gorm:"<-:create"`
	UpdatedAt    time.Time
	ModelRegisterUserActions
}

type CalendarPeriods struct {
	ModelEntity
	CalendarID   *string                    `gorm:"column:calendar_id" gorm:"foreignKey:CalendarId;references:ID"`
	Calendar     Calendars                  `gorm:"calendar_id" gorm:"index"`
	PeriodNumber CalendarPeriodPeriodNumber `gorm:"period_number"`
	Description  string                     `gorm:"description"`
	Year         int                        `gorm:"year"`
	DayType      CalendarPeriodDayType      `gorm:"day_type"`
	StartHour    int                        `gorm:"start_hour"`
	EndHour      int                        `gorm:"end_hour"`
	Energy       bool                       `gorm:"energy"`
	Power        bool                       `gorm:"power"`
	StartDate    time.Time                  `gorm:"start_date"`
	EndDate      time.Time                  `gorm:"end_date"`
	ModelRegisterUserActions
}
type Tariff struct {
	ID           string              `gorm:"tariff_id" gorm:"primaryKey" gorm:"index"`
	Description  string              `gorm:"description"`
	TensionLevel TariffsTensionLevel `gorm:"tension_level"`
	CodeOdos     string              `gorm:"code_odos"`
	CodeOne      string              `gorm:"code_one"`
	Periods      int                 `gorm:"periods"`
	Coef         TariffsCoef         `gorm:"column:coef"`
	GeographicID *string             `gorm:"column:geographic_id;type:uuid"  gorm:"foreignKey:GeographicId;references:ID"`
	Geographic   GeographicZones     `gorm:"geographic_id"`
	CalendarID   *string             `gorm:"column:calendar_id" gorm:"foreignKey:CalendarId;references:ID"`
	Calendar     Calendars           `gorm:"calendar_id" gorm:"index"`
	SoftDelete
	ModelRegisterUserActions
}
type TariffCalendar struct {
	ModelEntity
	CalendarID   *string         `gorm:"column:calendar_id" gorm:"foreignKey:CalendarId;references:ID"`
	Calendar     Calendars       `gorm:"calendar_id" gorm:"index"`
	TariffIdID   *string         `gorm:"column:tariff_id" gorm:"foreignKey:TariffId;references:ID"`
	TariffId     Tariff          `gorm:"tariff_id" gorm:"index"`
	StartDate    time.Time       `gorm:"start_date"`
	EndDate      time.Time       `gorm:"end_date"`
	GeographicID *string         `gorm:"column:geographic_id;type:uuid"  gorm:"foreignKey:GeographicId;references:ID"`
	Geographic   GeographicZones `gorm:"geographic_id"`
	ModelRegisterUserActions
}

type Seasons struct {
	ModelEntity
	Name         string          `gorm:"name"`
	Description  string          `gorm:"description"`
	GeographicID *string         `gorm:"column:geographic_id;type:uuid"  gorm:"foreignKey:GeographicId;references:ID"`
	Geographic   GeographicZones `gorm:"geographic_id"`
	ModelRegisterUserActions
}

type DayTypes struct {
	ModelEntity
	SeasonsID *string `gorm:"column:seasons_id;type:uuid"  gorm:"foreignKey:SeasonId;references:ID"`
	Seasons   Seasons `gorm:"seasons_id"`
	Name      string  `gorm:"name"`
	Month     int     `gorm:"month"`
	IsFestive bool    `gorm:"column:is_festive"`
	ModelRegisterUserActions
}
type FestiveDays struct {
	ModelEntity
	Date         time.Time       `gorm:"column:date"`
	Description  string          `gorm:"column:description"`
	GeographicID *string         `gorm:"column:geographic_id;type:uuid"  gorm:"foreignKey:GeographicId;references:ID"`
	Geographic   GeographicZones `gorm:"geographic_id"`
	ModelRegisterUserActions
}

// Defines values for CalendarPeriodDayType.
const (
	Festive  CalendarPeriodDayType = "Festive"
	Workable CalendarPeriodDayType = "Workable"
)

// Defines values for CalendarPeriodPeriodNumber.
const (
	P1 CalendarPeriodPeriodNumber = "P1"
	P2 CalendarPeriodPeriodNumber = "P2"
	P3 CalendarPeriodPeriodNumber = "P3"
	P4 CalendarPeriodPeriodNumber = "P4"
	P5 CalendarPeriodPeriodNumber = "P5"
	P6 CalendarPeriodPeriodNumber = "P6"
)

// Defines values for TariffsTensionLevel.
const (
	TariffsTensionLevelAT TariffsTensionLevel = "AT"

	TariffsTensionLevelBT TariffsTensionLevel = "BT"

	TariffsTensionLevelMT TariffsTensionLevel = "MT"
)

// Defines values for TariffsCoef.
const (
	TariffsCoefA TariffsCoef = "A"

	TariffsCoefB TariffsCoef = "B"

	TariffsCoefC TariffsCoef = "C"

	TariffsCoefD TariffsCoef = "D"
)

// CalendarPeriodDayType defines model for CalendarPeriod.DayType.
type CalendarPeriodDayType string

// CalendarPeriodPeriodNumber defines model for CalendarPeriod.PeriodNumber.
type CalendarPeriodPeriodNumber string

// TariffsTensionLevel defines model for tariffs.TensionLevel.
type TariffsTensionLevel string

// TariffsCoef defines model for Tariffs.Coef.
type TariffsCoef string
