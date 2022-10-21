package billing_measures

import (
	"context"
	"time"
)

type QueryConsumCoefficient struct {
	EndDate   time.Time
	StartDate time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ConsumCoefficientRepository
type ConsumCoefficientRepository interface {
	Search(ctx context.Context, q QueryConsumCoefficient) (float64, error)
}
