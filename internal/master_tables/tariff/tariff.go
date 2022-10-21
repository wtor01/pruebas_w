package tariff

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"time"
)

type Search struct {
	Q           string
	Limit       int
	Offset      *int
	Sort        map[string]string
	CurrentUser auth.User
}

type Tariffs struct {
	Id           string
	Description  string
	TensionLevel string
	CodeOdos     string
	CodeOne      string
	Periods      int
	GeographicId string
	CalendarId   string
	Coef         string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedByID  *string
	UpdatedBy    string
	UpdatedAt    time.Time
}
type TariffsCalendar struct {
	CalendarId     string
	TariffId       string
	StartDate      time.Time
	EndDate        time.Time
	GeoGraphicCode string
}
