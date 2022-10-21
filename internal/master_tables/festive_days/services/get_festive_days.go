package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"context"
)

type GetListFestiveDaysDto struct {
	Q             string
	Limit         int
	Offset        *int
	Sort          map[string]string
	DistributorID string
}

type GetListFestiveDays struct {
	repository festive_days.FestiveDayRepository
}

func NewGetListFestiveDaysService(repository festive_days.FestiveDayRepository) *GetListFestiveDays {
	return &GetListFestiveDays{repository: repository}
}

func (s GetListFestiveDays) Handler(ctx context.Context, dto GetListFestiveDaysDto) ([]festive_days.FestiveDay, int, error) {

	result, count, err := s.repository.GetListFestiveDays(ctx, festive_days.Search{
		Q:      dto.Q,
		Limit:  dto.Limit,
		Offset: dto.Offset,
		Sort:   dto.Sort,
	})

	if err != nil {
		return []festive_days.FestiveDay{}, 0, err
	}

	return result, count, err
}
