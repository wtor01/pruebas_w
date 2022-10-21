package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"time"
)

func NewConsumProfile(
	date time.Time,
	version int,
	t billing_measures.ConsumProfileType,
	coefA float64,
	coefB float64,
	coefC float64,
	coefD float64,
) billing_measures.ConsumProfile {
	c := billing_measures.ConsumProfile{Date: date, Version: version, Type: t, CoefA: &coefA, CoefB: &coefB, CoefC: &coefC, CoefD: &coefD}
	if coefA == -1 {
		c.CoefA = nil
	}
	if coefB == -1 {
		c.CoefB = nil
	}
	if coefC == -1 {
		c.CoefC = nil
	}
	if coefD == -1 {
		c.CoefD = nil
	}
	return c
}

var ConsumProfileInsertFile = []billing_measures.ConsumProfile{
	NewConsumProfile(
		time.Date(2022, 01, 01, 00, 0, 0, 0, time.UTC),
		1,
		"PROVISIONAL",
		0.000135829136,
		0.000087585352,
		0.000051169015,
		-1,
	),
	NewConsumProfile(
		time.Date(2022, 01, 01, 00, 0, 0, 0, time.UTC),
		1,
		"INITIAL",
		0.000135829136,
		0.000087585352,
		0.000051169015,
		0.000051169015,
	),
}
