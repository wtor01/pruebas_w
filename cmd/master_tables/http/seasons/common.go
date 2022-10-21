package seasons

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
)

func SeasonsToResponse(d seasons.Seasons) Seasons {
	return Seasons{
		Id:             d.ID.String(),
		Name:           d.Name,
		Description:    d.Description,
		GeographicCode: d.GeographicCode,
	}
}

func DayTypesToResponse(season seasons.DayTypes) DayTypes {
	return DayTypes{
		Id:        season.ID,
		Name:      season.Name,
		IsFestive: season.IsFestive,
		Month:     season.Month,
		SeasonId:  season.SeasonsId,
	}
}
