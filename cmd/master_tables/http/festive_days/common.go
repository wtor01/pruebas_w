package festive_days

import "bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"

func FestiveDaysToResponse(festiveDay festive_days.FestiveDay) FestiveDays {
	return FestiveDays{
		Id:           festiveDay.Id.String(),
		Date:         festiveDay.Date.String(),
		Description:  festiveDay.Description,
		GeographicId: festiveDay.GeographicId,
	}
}
