package process_measures

import (
	"context"
	"time"
)

type GetClosure struct {
	Id            string
	DistributorId string
	CUPS          string
	StartDate     time.Time
	EndDate       time.Time
	Moment        SelectMoment
}
type GetResume struct {
	Cups          string
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

type SelectMoment string

const (
	Next   SelectMoment = "next"
	Before SelectMoment = "before"
	Actual SelectMoment = "actual"
)

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ProcessMeasureClosureRepository

type ProcessMeasureClosureRepository interface {
	GetClosure(ctx context.Context, query GetClosure) (ProcessedMonthlyClosure, error)
	GetResume(ctx context.Context, query GetResume) (ResumesProcessMonthlyClosure, error)
	CreateClosure(ctx context.Context, monthly ProcessedMonthlyClosure) error
	UpdateClosure(ctx context.Context, id string, monthly ProcessedMonthlyClosure) error
	GetClosureOne(ctx context.Context, id string) (ProcessedMonthlyClosure, error)
}
