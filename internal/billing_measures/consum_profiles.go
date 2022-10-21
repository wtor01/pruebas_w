package billing_measures

import (
	"context"
	"time"
)

type ConsumProfileType string

const (
	InitialTypeConsumProfile     ConsumProfileType = "INITIAL"
	ProvisionalTypeConsumProfile ConsumProfileType = "PROVISIONAL"
	FinalTypeConsumProfile       ConsumProfileType = "FINAL"
)

type ConsumProfile struct {
	Date    time.Time         `json:"date" bson:"date"`
	Version int               `json:"version" bson:"version"`
	Type    ConsumProfileType `json:"type" bson:"type"`
	CoefA   *float64          `json:"coef_a" bson:"coef_a"`
	CoefB   *float64          `json:"coef_b" bson:"coef_b"`
	CoefC   *float64          `json:"coef_c" bson:"coef_c"`
	CoefD   *float64          `json:"coef_d" bson:"coef_d"`
}

func (c ConsumProfile) GetValueByCoefficientType(t CoefficientType) float64 {
	var value *float64
	switch t {
	case CoefficientA:
		value = c.CoefA
	case CoefficientB:
		value = c.CoefB
	case CoefficientC:
		value = c.CoefC
	case CoefficientD:
		value = c.CoefD
	}
	if value == nil {
		return 0.0
	}
	return *value
}

type QueryConsumProfile struct {
	EndDate   time.Time
	StartDate time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ConsumProfileRepository
type ConsumProfileRepository interface {
	Save(ctx context.Context, profile ConsumProfile) error
	Search(ctx context.Context, q QueryConsumProfile) ([]ConsumProfile, error)
}
