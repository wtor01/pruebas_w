package clients

import (
	"context"
)

type GetTariffDto struct {
	ID string
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
}

//go:generate mockery --case=snake --outpkg=mocks --output=./clients_mocks --name=MasterTables
type MasterTables interface {
	GetTariff(ctx context.Context, dto GetTariffDto) (Tariffs, error)
}
