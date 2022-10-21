package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sort"
	"strings"
	"time"
)

type SearchFiscalBillingMeasuresDashboardDTO struct {
	Id            *string
	Cups          string
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

type SearchFiscalBillingMeasures struct {
	Location   *time.Location
	repository billing_measures.BillingMeasuresDashboardRepository
	tracer     trace.Tracer
}

func NewSearchFiscalBillingMeasuresDashboard(repository billing_measures.BillingMeasuresDashboardRepository, location *time.Location) *SearchFiscalBillingMeasures {
	return &SearchFiscalBillingMeasures{
		repository: repository,
		Location:   location,
		tracer:     telemetry.GetTracer(),
	}
}

func (s SearchFiscalBillingMeasures) Handler(ctx context.Context, dto SearchFiscalBillingMeasuresDashboardDTO) ([]billing_measures.FiscalBillingMeasuresDashboard, error) {
	ctx, span := s.tracer.Start(ctx, "SearchFiscalBillingMeasures - Handler")
	defer span.End()

	billingMeasures, err := s.repository.SearchFiscalBillingMeasures(ctx, dto.Cups, dto.DistributorId, dto.StartDate, dto.EndDate)
	if err != nil {
		return []billing_measures.FiscalBillingMeasuresDashboard{}, err
	}
	lastBillingMeasures, err := s.repository.SearchLastBillingMeasures(ctx, dto.Cups, dto.DistributorId)
	if err != nil {
		return []billing_measures.FiscalBillingMeasuresDashboard{}, err
	}
	response := make([]billing_measures.FiscalBillingMeasuresDashboard, 0, len(billingMeasures))

	span.SetAttributes(attribute.Int("response", len(response)))

	for _, billingMeasure := range billingMeasures {
		fiscalBilling := s.toResponse(ctx, billingMeasure, lastBillingMeasures)
		response = append(response, fiscalBilling)
	}
	return response, nil
}

func (s SearchFiscalBillingMeasures) toResponse(ctx context.Context, billingMeasure, lastBillingMeasure billing_measures.BillingMeasure) billing_measures.FiscalBillingMeasuresDashboard {
	ctx, span := s.tracer.Start(ctx, "SearchFiscalBillingMeasures - toResponse")
	defer span.End()

	fiscalBilling := billing_measures.NewFiscalBillingMeasuresDashboard(billingMeasure, lastBillingMeasure)

	for _, period := range billingMeasure.GetPeriods() {

		billingBalancePeriod := billingMeasure.GetBalancePeriod(period)
		balancePeriod := billing_measures.BalancePeriod{
			AE:            billingBalancePeriod.AE,
			BalanceTypeAE: billingBalancePeriod.BalanceGeneralTypeAE,
			AI:            billingBalancePeriod.AI,
			BalanceTypeAI: billingBalancePeriod.BalanceGeneralTypeAI,
			R1:            billingBalancePeriod.R1,
			BalanceTypeR1: billingBalancePeriod.BalanceGeneralTypeR1,
			R2:            billingBalancePeriod.R2,
			BalanceTypeR2: billingBalancePeriod.BalanceGeneralTypeR2,
			R3:            billingBalancePeriod.R3,
			BalanceTypeR3: billingBalancePeriod.BalanceGeneralTypeR3,
			R4:            billingBalancePeriod.R4,
			BalanceTypeR4: billingBalancePeriod.BalanceGeneralTypeR4,
		}
		fiscalBilling.SetBalancePeriod(period, balancePeriod)
	}

	calendarCurveMap := make(map[string]billing_measures.GeneralEstimateMethod)
	curveGraphMap := make(map[string][]billing_measures.Values)

	for i, curve := range billingMeasure.BillingLoadCurve {
		dayToSearch, hourlyDay := s.getFormattedDate(curve.EndDate)

		curveMethod, _ := billingMeasure.GetLoadCurveGeneralEstimatedMethod(i, fiscalBilling.PrincipalMagnitude)

		fiscalBilling.AddConsum(curve.Period, billingMeasure.GetLoadCurvePeriodMagnitude(i, fiscalBilling.PrincipalMagnitude))
		fiscalBilling.AddSummary(curve.Period, curveMethod)
		s.setCalendarStatus(curveMethod, dayToSearch, calendarCurveMap)
		s.setCurveGraphMap(curve, curveMethod, dayToSearch, hourlyDay, curveGraphMap)
	}

	fiscalBilling.CalcTotalConsum(billingMeasure.Periods)

	curvesGraph := s.transformCurveGraph(curveGraphMap)
	fiscalBilling.SetCurveGraph(curvesGraph)

	calendarCurves := s.transformCalendarCurve(calendarCurveMap)
	fiscalBilling.SetCalendarCurve(calendarCurves)

	fiscalBilling.CalculatePercentages(fiscalBilling.Periods)

	return fiscalBilling
}

func (s SearchFiscalBillingMeasures) getFormattedDate(date time.Time) (string, string) {
	dateFormat := date.In(s.Location).Format("2006-01-02 15:04")
	dateFormatSplited := strings.Split(dateFormat, " ")
	return dateFormatSplited[0], dateFormatSplited[1]
}

func (s SearchFiscalBillingMeasures) setCurveGraphMap(curve billing_measures.BillingLoadCurve, status billing_measures.GeneralEstimateMethod, dayToSearch, hourlyDate string, curveGraphMap map[string][]billing_measures.Values) {
	_, ok := curveGraphMap[dayToSearch]
	if !ok {
		curveValues := make([]billing_measures.Values, 0, 24)
		curveValues = append(curveValues, billing_measures.Values{
			Hour:   hourlyDate,
			Status: status,
			AI:     curve.AI,
			AE:     curve.AE,
			AiAuto: curve.AIAuto,
			AeAuto: curve.AeAuto,
			R1:     curve.R1,
			R2:     curve.R2,
			R3:     curve.R3,
			R4:     curve.R4,
		})
		curveGraphMap[dayToSearch] = curveValues
		return

	}
	curveGraphMap[dayToSearch] = append(curveGraphMap[dayToSearch], billing_measures.Values{
		Hour:   hourlyDate,
		Status: status,
		AI:     curve.AI,
		AE:     curve.AE,
		R1:     curve.R1,
		R2:     curve.R2,
		R3:     curve.R3,
		R4:     curve.R4,
	})

}

func (s SearchFiscalBillingMeasures) setCalendarStatus(status billing_measures.GeneralEstimateMethod, dayToSearch string, calendarDays map[string]billing_measures.GeneralEstimateMethod) {
	calCurve, ok := calendarDays[dayToSearch]
	if !ok {
		calendarDays[dayToSearch] = status
		return
	}

	if calCurve == billing_measures.GeneralEstimated || status == "" {
		return
	}

	if !calCurve.Compare(billing_measures.GeneralMethodValue, status) || calCurve == "" {
		calendarDays[dayToSearch] = status
	}
}

func (s SearchFiscalBillingMeasures) transformCalendarCurve(curveMap map[string]billing_measures.GeneralEstimateMethod) []billing_measures.CalendarCurve {
	calendarCurve := make([]billing_measures.CalendarCurve, 0, len(curveMap))
	for date, status := range curveMap {
		calendarCurve = append(calendarCurve, billing_measures.CalendarCurve{date, status})
	}
	sort.Slice(calendarCurve, func(i, j int) bool {
		return calendarCurve[i].Date < calendarCurve[j].Date
	})
	return calendarCurve
}

func (s SearchFiscalBillingMeasures) transformCurveGraph(curveMap map[string][]billing_measures.Values) []billing_measures.Curve {
	curveGraph := make([]billing_measures.Curve, 0, len(curveMap))
	for date, values := range curveMap {
		curveGraph = append(curveGraph, billing_measures.Curve{Date: date, Values: values})
	}
	sort.Slice(curveGraph, func(i, j int) bool {
		return curveGraph[i].Date < curveGraph[j].Date
	})

	return curveGraph
}
