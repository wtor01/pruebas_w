package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"context"
)

type GetFestiveDayByIdDto struct {
	Id string
}

type GetFestiveDayById struct {
	repository festive_days.FestiveDayRepository
}

func NewGetFestiveDayByIdService(repository festive_days.FestiveDayRepository) *GetFestiveDayById {
	return &GetFestiveDayById{repository: repository}
}

func (s GetFestiveDayById) Handler(ctx context.Context, dto GetFestiveDayByIdDto) (festive_days.FestiveDay, error) {
	ds, err := s.repository.GetFestiveDayById(ctx, dto.Id)

	return ds, err
}
