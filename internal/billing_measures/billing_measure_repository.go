package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"time"
)

type QueryLast struct {
	CUPS string
	Date time.Time
}

type GetPrevious struct {
	CUPS     string
	InitDate time.Time
	EndDate  time.Time
}

type QueryLastHistory struct {
	CUPS       string
	InitDate   time.Time
	EndDate    time.Time
	Periods    []measures.PeriodKey
	Magnitudes []measures.Magnitude
}

type QueryGetCloseHistories struct {
	CUPS       string
	Periods    []measures.PeriodKey
	EndDate    time.Time
	Magnitudes []measures.Magnitude
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=BillingMeasureRepository
type BillingMeasureRepository interface {
	Last(ctx context.Context, q QueryLast) (BillingMeasure, error)
	LastHistory(ctx context.Context, q QueryLastHistory) (BillingMeasure, error)
	GetCloseHistories(ctx context.Context, q QueryGetCloseHistories) ([]BillingMeasure, error)
	Find(ctx context.Context, id string) (BillingMeasure, error)
	Save(ctx context.Context, measure BillingMeasure) error
	SaveAll(ctx context.Context, measures []BillingMeasure) error
	GetPrevious(ctx context.Context, query GetPrevious) (BillingMeasure, error)
}
