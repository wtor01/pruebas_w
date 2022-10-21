package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type SearchServicePointProcessMeasuresDashboardCurvesDTO struct {
	DistributorId string
	Cups          string
	StartDate     time.Time
	EndDate       time.Time
	CurveType     measures.MeasureCurveReadingType
}

type SearchServicePointDashboardCurves struct {
	repositoryProcessMeasures process_measures.ProcessedMeasureRepository
	tracer                    trace.Tracer
	Location                  *time.Location
}

func NewSearchServicePointProcessMeasuresDashboardCurvesDTO(distributorId, cups string, startDate, endDate time.Time, registerType measures.MeasureCurveReadingType) SearchServicePointProcessMeasuresDashboardCurvesDTO {
	return SearchServicePointProcessMeasuresDashboardCurvesDTO{
		DistributorId: distributorId,
		Cups:          cups,
		StartDate:     startDate,
		EndDate:       endDate,
		CurveType:     registerType,
	}
}

func NewSearchServicePointDashboardCurves(repositoryProcessMeasures process_measures.ProcessedMeasureRepository, location *time.Location) *SearchServicePointDashboardCurves {
	return &SearchServicePointDashboardCurves{repositoryProcessMeasures: repositoryProcessMeasures, Location: location, tracer: telemetry.GetTracer()}
}

func (s SearchServicePointDashboardCurves) Handler(ctx context.Context, dto SearchServicePointProcessMeasuresDashboardCurvesDTO) ([]process_measures.ServicePointDashboardCurve, error) {
	ctx, span := s.tracer.Start(ctx, "SearchServicePointDashboardCurves - Handler")
	defer span.End()

	curves, err := s.repositoryProcessMeasures.ProcessedLoadCurveByCups(
		ctx,
		process_measures.QueryProcessedLoadCurveByCups{
			CUPS:      dto.Cups,
			StartDate: dto.StartDate,
			EndDate:   dto.EndDate,
			CurveType: dto.CurveType,
		})

	if err != nil {
		return []process_measures.ServicePointDashboardCurve{}, err
	}

	result := make([]process_measures.ServicePointDashboardCurve, 0)
	for _, curve := range curves {
		valuesCurve := measures.Values{
			AI: curve.AI,
			AE: curve.AE,
			R1: curve.R1,
			R2: curve.R2,
			R3: curve.R3,
			R4: curve.R4,
		}

		/*
		* - Si el origin es filled tenemos que poner su status a NONE
		* - Si el origin no es filled ponemos el status que venga por defecto
		 */
		spdc := process_measures.ServicePointDashboardCurve{
			Date: curve.EndDate,
			Values: struct {
				measures.Values
			}{valuesCurve},
		}
		if curve.Origin == measures.Filled {
			spdc.Status = measures.None
		} else {
			spdc.Status = curve.ValidationStatus
		}

		result = append(result, spdc)

	}
	span.SetAttributes(attribute.Int("response", len(result)))

	return result, nil
}
