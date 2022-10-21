package measures

import (
	"context"
	"time"
)

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ReprocessingDateRepository
type ReprocessingDateRepository interface {
	GetDate(ctx context.Context, keySearch string) (error, time.Time)
	SetDate(ctx context.Context, keySearch string, date time.Time) error
}
