package redis

import (
	"context"
	"time"
)

type ReprocessingDate struct {
	r *DataCacheRedis
}

func NewDateRedis(r *DataCacheRedis) ReprocessingDate {
	return ReprocessingDate{r: r}
}
func (r ReprocessingDate) GetDate(ctx context.Context, keySearch string) (error, time.Time) {
	date := time.Time{}
	err := r.r.Get(ctx, &date, keySearch)
	return err, date
}
func (r ReprocessingDate) SetDate(ctx context.Context, keySearch string, date time.Time) error {
	return r.r.Set(ctx, keySearch, date)
}
