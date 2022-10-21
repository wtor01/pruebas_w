package tariff

import (
	"context"
)

type RepositoryTariff interface {
	GetAllTariffs(ctx context.Context, search Search) ([]Tariffs, int, error)
	InsertTariff(ctx context.Context, tf Tariffs) error
	DeleteTariff(ctx context.Context, tariffId string) error
	GetTariff(ctx context.Context, tariffId string) (Tariffs, error)
	ModifyTariff(ctx context.Context, tariffId string, tf Tariffs) error
	InsertTariffCalendar(ctx context.Context, calendar_id string, tariff Tariffs, created_by string) error
	GetAllTariffCalendar(ctx context.Context, code string, search Search) ([]TariffsCalendar, int, error)
}
