package smarkia

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"context"
	"time"
)

type Distributor struct {
	SmarkiaCode string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=GetDistributorser
type GetDistributorser interface {
	GetDistributors(ctx context.Context) ([]Distributor, error)
}

type Ct struct {
	ID string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=GetCTer
type GetCTer interface {
	GetCTs(ctx context.Context, distributorID string) ([]Ct, error)
}

type Equipment struct {
	ID   string
	CUPS string
	CtId string
}

type GetEquipmentsQuery struct {
	Id   string
	Cups string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=GetEquipmentser
type GetEquipmentser interface {
	GetEquipments(ctx context.Context, query GetEquipmentsQuery) ([]Equipment, error)
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=GetMagnitudeser
type GetMagnitudeser interface {
	GetMagnitudes(ctx context.Context, equipmentID string, date time.Time) ([]gross_measures.MeasureCurveWrite, error)
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=GetClosinger
type GetClosinger interface {
	GetClosinger(ctx context.Context, equipmentID string, date time.Time) ([]gross_measures.MeasureCloseWrite, error)
}
