package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/trace"
)

type FestiveDaysGenerateService struct {
	festiveDaysRepository    festive_days.FestiveDayRepository
	calendarPeriodRepository measures.CalendarPeriodRepository
	tracer                   trace.Tracer
}

func NewFestiveDaysGenerateService(festiveDaysRepository festive_days.FestiveDayRepository, calendarPeriodRepository measures.CalendarPeriodRepository) *FestiveDaysGenerateService {
	return &FestiveDaysGenerateService{
		festiveDaysRepository:    festiveDaysRepository,
		calendarPeriodRepository: calendarPeriodRepository,
		tracer:                   telemetry.GetTracer(),
	}
}

func (s FestiveDaysGenerateService) Handler(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "FestiveDaysGenerate - Handler")
	defer span.End()

	festiveDays, _, err := s.festiveDaysRepository.GetListFestiveDays(ctx, festive_days.Search{
		Q:     "",
		Limit: 365,
	})

	if err != nil {
		return err
	}

	err = s.calendarPeriodRepository.DeleteFestiveDays(ctx)

	if err != nil {
		return err
	}

	for _, festive := range festiveDays {
		if festive.GeographicId == "" {
			continue
		}

		err = s.calendarPeriodRepository.SaveFestiveDay(ctx, measures.FestiveDay{
			Date:         festive.Date,
			GeographicID: festive.GeographicId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
