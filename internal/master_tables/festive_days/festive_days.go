package festive_days

import (
	"context"
	"strings"
)

type Search struct {
	Q      string
	Limit  int
	Offset *int
	Sort   map[string]string
	date   string
}

func NewSort(arr []string) map[string]string {
	sort := make(map[string]string)

	for _, param := range arr {
		sortParam := strings.Split(param, " ")

		if len(sortParam) != 2 {
			continue
		}

		if sortParam[1] != "asc" && sortParam[1] != "desc" {
			continue
		}

		sort[sortParam[0]] = sortParam[1]
	}
	return sort
}

type FestiveDayRepository interface {
	GetListFestiveDays(ctx context.Context, search Search) ([]FestiveDay, int, error)
	GetFestiveDayById(ctx context.Context, festiveDayId string) (FestiveDay, error)
	SaveFestiveDay(ctx context.Context, festiveDay FestiveDay) error
	UpdateFestiveDay(ctx context.Context, festiveDayId string, festiveDay FestiveDay) error
	DeleteFestiveDay(ctx context.Context, festiveDayId string) error
}
