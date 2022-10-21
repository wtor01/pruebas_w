package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
)

type IsCloseMeasureComplete struct {
	ID        string
	b         *BillingMeasure    `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
}

func NewIsCloseMeasureComplete(b *BillingMeasure, magnitude measures.Magnitude) *IsCloseMeasureComplete {
	return &IsCloseMeasureComplete{
		b:         b,
		ID:        "IS_CLOSE_MEASURE_COMPLETE",
		Magnitude: magnitude,
	}
}

func (i IsCloseMeasureComplete) Eval(_ context.Context) bool {

	for _, p := range i.b.Periods {
		if i.b.GetBalanceStatus(p, i.Magnitude) == measures.Invalid {
			return false
		}
	}
	return true
}

type IsMoreThanOneMissingPeriod struct {
	ID        string
	b         *BillingMeasure    `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
}

func NewIsMoreThanOneMissingPeriod(b *BillingMeasure, magnitude measures.Magnitude) *IsMoreThanOneMissingPeriod {
	return &IsMoreThanOneMissingPeriod{
		b:         b,
		ID:        "IS_CLOSE_MEASURE_COMPLETE",
		Magnitude: magnitude,
	}
}

func (i IsMoreThanOneMissingPeriod) Eval(_ context.Context) bool {
	count := 0
	for _, p := range i.b.Periods {
		if i.b.GetBalanceStatus(p, i.Magnitude) == measures.Invalid {
			count++
		}
		if count >= 2 {
			return true
		}
	}
	return false
}

type IsHourly struct {
	ID string
	b  *BillingMeasure `bson:"-"`
}

func NewIsHourly(b *BillingMeasure) *IsHourly {
	return &IsHourly{
		ID: "IS_HOURLY",
		b:  b,
	}
}

func (i IsHourly) Eval(_ context.Context) bool {
	return i.b.IsHourly()
}

type HasClosedHistory struct {
	ID                string
	b                 *BillingMeasure          `bson:"-"`
	billingRepository BillingMeasureRepository `bson:"-"`
	Magnitude         measures.Magnitude       `bson:"-"`
	ContextNoTlg      *GraphContext            `bson:"-"`
}

func NewHasClosedHistory(
	b *BillingMeasure,
	repository BillingMeasureRepository,
	magnitude measures.Magnitude,
	contextNoTlg *GraphContext,
) *HasClosedHistory {
	return &HasClosedHistory{
		ID:                "HAS_CLOSED_HISTORY",
		b:                 b,
		billingRepository: repository,
		Magnitude:         magnitude,
		ContextNoTlg:      contextNoTlg,
	}
}

func (i HasClosedHistory) Eval(ctx context.Context) bool {
	if !i.ContextNoTlg.IsClosedHistoryRequested {
		billingHistories, err := i.billingRepository.GetCloseHistories(ctx, QueryGetCloseHistories{
			CUPS:       i.b.CUPS,
			Periods:    i.b.Periods,
			EndDate:    i.b.EndDate,
			Magnitudes: []measures.Magnitude{i.Magnitude},
		})

		if err == nil {
			i.ContextNoTlg.ClosedHistory = billingHistories
		}

		i.ContextNoTlg.IsClosedHistoryRequested = true
	}

	return len(i.ContextNoTlg.ClosedHistory) >= 4
}

type HasAnualHistory struct {
	ID                       string
	b                        *BillingMeasure          `bson:"-"`
	billingMeasureRepository BillingMeasureRepository `bson:"-"`
	Magnitude                measures.Magnitude       `bson:"-"`
	ContextNoTlg             *GraphContext            `bson:"-"`
}

func NewHasAnualHistory(
	b *BillingMeasure,
	repository BillingMeasureRepository,
	magnitude measures.Magnitude,
	contextNoTlg *GraphContext,
) *HasAnualHistory {
	return &HasAnualHistory{
		ID:                       "HAS_ANUAL_HISTORY",
		b:                        b,
		billingMeasureRepository: repository,
		Magnitude:                magnitude,
		ContextNoTlg:             contextNoTlg,
	}
}

func (i HasAnualHistory) Eval(ctx context.Context) bool {
	if !i.ContextNoTlg.IsAnualHistoryRequested {
		billingHistory, err := i.billingMeasureRepository.LastHistory(ctx, QueryLastHistory{
			CUPS:       i.b.CUPS,
			InitDate:   i.b.EndDate,
			EndDate:    i.b.EndDate,
			Periods:    i.b.Periods,
			Magnitudes: []measures.Magnitude{i.Magnitude},
		})

		if err == nil {
			i.ContextNoTlg.AnualHistory = billingHistory
		}

		i.ContextNoTlg.IsAnualHistoryRequested = true
	}

	return i.ContextNoTlg.AnualHistory.Id != ""
}

type IsHouseCloseAndCloseAtrAreEmpty struct {
	ID        string
	b         *BillingMeasure    `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
}

func NewIsHouseCloseAndCloseAtrAreEmpty(b *BillingMeasure, magnitude measures.Magnitude) *IsHouseCloseAndCloseAtrAreEmpty {
	return &IsHouseCloseAndCloseAtrAreEmpty{
		ID:        "IS_HOUSE_CLOSE_AND_CLOSE_ATR_Are_EMPTY",
		b:         b,
		Magnitude: magnitude,
	}
}

func (i IsHouseCloseAndCloseAtrAreEmpty) Eval(_ context.Context) bool {
	return i.b.isPointType([]measures.PointType{"3", "4"}) && i.b.IsHouseClose() && i.b.AreAllCloseAtrEmpty(i.Magnitude)
}

type HasIterativeHistory struct {
	ID                         string
	b                          *BillingMeasure                             `bson:"-"`
	processedMeasureRepository process_measures.ProcessedMeasureRepository `bson:"-"`
	minCount                   int                                         `bson:"-"`
	Magnitude                  measures.Magnitude                          `bson:"-"`
	ContextNoTlg               *GraphContext                               `bson:"-"`
}

func NewHasIterativeHistory(
	b *BillingMeasure,
	processedMeasureRepository process_measures.ProcessedMeasureRepository,
	magnitude measures.Magnitude,
	contextNoTlg *GraphContext,
) *HasIterativeHistory {
	return &HasIterativeHistory{
		b:                          b,
		processedMeasureRepository: processedMeasureRepository,
		ID:                         "HAS_ITERATIVE_HISTORY",
		minCount:                   6,
		Magnitude:                  magnitude,
		ContextNoTlg:               contextNoTlg,
	}
}

func (i HasIterativeHistory) Eval(ctx context.Context) bool {

	if !i.ContextNoTlg.IsIterativeHistoryRequested {
		var loadCurves []process_measures.ProcessedLoadCurve

		previousLoadCurve, err := i.processedMeasureRepository.ListHistoryLoadCurve(ctx, process_measures.QueryHistoryLoadCurve{
			CUPS:          i.b.CUPS,
			StartDate:     i.b.InitDate.AddDate(0, -6, 0),
			EndDate:       i.b.InitDate,
			Periods:       i.b.Periods,
			WithCriterias: true,
			Magnitude:     i.Magnitude,
		})

		if err == nil {
			loadCurves = previousLoadCurve
		}

		nextLoadCurve, err := i.processedMeasureRepository.ListHistoryLoadCurve(ctx, process_measures.QueryHistoryLoadCurve{
			CUPS:          i.b.CUPS,
			StartDate:     i.b.EndDate,
			EndDate:       i.b.EndDate.AddDate(0, 6, 0),
			Periods:       i.b.Periods,
			WithCriterias: true,
			IsFuture:      true,
			Magnitude:     i.Magnitude,
		})

		if err == nil {
			loadCurves = append(loadCurves, nextLoadCurve...)
		}

		i.ContextNoTlg.IterativeHistory = loadCurves
		i.ContextNoTlg.IsIterativeHistoryRequested = true
	}

	return len(i.ContextNoTlg.IterativeHistory) >= (i.minCount * len(i.b.Periods))
}

type IsSimpleHistoric struct {
	ID                         string
	Magnitude                  measures.Magnitude                          `bson:"-"`
	B                          *BillingMeasure                             `bson:"-"`
	ContextNoTlg               *GraphContext                               `bson:"-"`
	processedMeasureRepository process_measures.ProcessedMeasureRepository `bson:"-"`
}

func NewIsSimpleHistoric(
	B *BillingMeasure,
	contextNoTlg *GraphContext,
	magnitude measures.Magnitude,
	processedMeasureRepository process_measures.ProcessedMeasureRepository,
) *IsSimpleHistoric {
	return &IsSimpleHistoric{
		B:                          B,
		ContextNoTlg:               contextNoTlg,
		ID:                         "CHH_SIMPLE_HISTORIC",
		processedMeasureRepository: processedMeasureRepository,
		Magnitude:                  magnitude,
	}
}

func (conditional IsSimpleHistoric) Eval(ctx context.Context) bool {

	periods := conditional.B.Periods

	idxLastMeasurePeriod := map[measures.PeriodKey]int{}
	idxFirstMeasurePeriod := map[measures.PeriodKey]int{}
	if !conditional.ContextNoTlg.IsSimpleHistoricRequested {
		previousLoadCurve, err := conditional.processedMeasureRepository.ListHistoryLoadCurve(ctx, process_measures.QueryHistoryLoadCurve{
			CUPS:      conditional.B.CUPS,
			StartDate: conditional.B.InitDate.AddDate(0, 0, -1),
			EndDate:   conditional.B.InitDate,
			Count:     1,
			Periods:   conditional.B.GetPeriods(),
			Magnitude: conditional.Magnitude,
			IsFuture:  false,
		})

		if err == nil {
			conditional.ContextNoTlg.SimpleHistoric.PreviousLoadCurve = previousLoadCurve
		}
		nextLoadCurve, err := conditional.processedMeasureRepository.ListHistoryLoadCurve(ctx, process_measures.QueryHistoryLoadCurve{
			CUPS:      conditional.B.CUPS,
			StartDate: conditional.B.EndDate,
			EndDate:   conditional.B.EndDate.AddDate(0, 0, 1),
			Count:     1,
			Periods:   conditional.B.GetPeriods(),
			Magnitude: conditional.Magnitude,
			IsFuture:  true,
		})

		if err == nil {
			conditional.ContextNoTlg.SimpleHistoric.NextLoadCurve = nextLoadCurve
		}
		conditional.ContextNoTlg.IsSimpleHistoricRequested = true
	}

	for _, p := range periods {
		idxFirstMeasurePeriod[p] = -1
		idxLastMeasurePeriod[p] = -1
	}

	for p := range idxLastMeasurePeriod {
		for idx := range conditional.B.BillingLoadCurve {
			measure := conditional.B.BillingLoadCurve[len(conditional.B.BillingLoadCurve)-idx-1]
			if p != measure.Period {
				continue
			}

			if idxLastMeasurePeriod[p] == -1 {
				idxLastMeasurePeriod[p] = len(conditional.B.BillingLoadCurve) - idx - 1
			}
			continueBool := false

			for _, p := range periods {
				if idxLastMeasurePeriod[p] == -1 {
					continueBool = true
					break
				}
			}
			if continueBool == false {
				break
			}
		}
	}

	for p := range idxFirstMeasurePeriod {
		for idx, curve := range conditional.B.BillingLoadCurve {

			if p != curve.Period {
				continue
			}

			if idxFirstMeasurePeriod[p] == -1 {
				idxFirstMeasurePeriod[p] = idx
			}
			continueBool := false

			for _, pp := range periods {
				if idxFirstMeasurePeriod[pp] == -1 {
					continueBool = true
					break
				}
			}
			if continueBool == false {
				break
			}
		}
	}

	// SE COMPRUEBA SI POR CADA PRIMER VALOR POR CADA PERIODO TIENE UN VALOR ANTERIOR REAL
	for p := range idxFirstMeasurePeriod {
		if idxFirstMeasurePeriod[p] == -1 {
			continue
		}
		curve := conditional.B.BillingLoadCurve[idxFirstMeasurePeriod[p]]

		if curve.Origin == measures.Filled {
			if len(conditional.ContextNoTlg.SimpleHistoric.PreviousLoadCurve) == 0 {
				return false
			}
			for idx := range conditional.ContextNoTlg.SimpleHistoric.PreviousLoadCurve {
				measurePrevious := conditional.ContextNoTlg.SimpleHistoric.PreviousLoadCurve[len(conditional.ContextNoTlg.SimpleHistoric.PreviousLoadCurve)-idx-1]
				if measurePrevious.Period != curve.Period {
					continue
				}
				if measurePrevious.Origin == measures.Filled {
					return false
				}
				break
			}
		}
	}
	// SE COMPRUEBA SI POR CADA ULTIMO VALOR POR CADA PERIODO TIENE UN VALOR POSTERIOR REAL
	for p := range idxLastMeasurePeriod {
		if idxLastMeasurePeriod[p] == -1 {
			continue
		}

		curve := conditional.B.BillingLoadCurve[idxLastMeasurePeriod[p]]
		if curve.Origin == measures.Filled {

			if len(conditional.ContextNoTlg.SimpleHistoric.NextLoadCurve) == 0 {
				return false
			}

			for _, measureNext := range conditional.ContextNoTlg.SimpleHistoric.NextLoadCurve {
				if measureNext.Period != curve.Period {
					continue
				}
				if measureNext.Origin == measures.Filled {
					return false
				}
				break
			}
		}

	}

	return true
}

type IsBalanceValid struct {
	ID        string
	b         *BillingMeasure    `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
}

func NewIsBalanceValid(b *BillingMeasure, magnitude measures.Magnitude) *IsBalanceValid {
	return &IsBalanceValid{
		b:         b,
		ID:        "IS_BALANCE_VALID",
		Magnitude: magnitude,
	}
}

func (i IsBalanceValid) Eval(_ context.Context) bool {
	return i.b.IsBalanceValid(i.Magnitude)
}

func NewIsChhCompleted(b *BillingMeasure) *IsChhComplete {
	return &IsChhComplete{
		ID: "IS_CHH_COMPLETE",
		b:  b,
	}
}
func (i IsChhComplete) Eval(_ context.Context) bool {
	return i.b.IsCchComplete()
}

type IsChhValid struct {
	ID string
	b  *BillingMeasure `bson:"-"`
}

func NewIsChhValid(b *BillingMeasure) *IsChhValid {
	return &IsChhValid{
		ID: "IS_CCH_VALID",
		b:  b,
	}
}
func (i IsChhValid) Eval(_ context.Context) bool {
	return i.b.IsChhValid()
}

type IsEmptyCentralHoursCch struct {
	ID string
	b  *BillingMeasure `bson:"-"`
}

func NewIsEmptyCentralHoursCch(b *BillingMeasure) *IsEmptyCentralHoursCch {
	return &IsEmptyCentralHoursCch{
		ID: "IS_EMPTY_CENTRAL_HOURS_CCH",
		b:  b,
	}
}
func (i IsEmptyCentralHoursCch) Eval(_ context.Context) bool {
	return i.b.IsEmptyCentralHoursCch()
}

type IsThereChhWindows struct {
	ID string
	b  *BillingMeasure `bson:"-"`
}

func NewIsThereChhWindows(b *BillingMeasure) *IsThereChhWindows {
	return &IsThereChhWindows{
		ID: "IS_THERE_CCH_WINDOWS",
		b:  b,
	}
}
func (i IsThereChhWindows) Eval(_ context.Context) bool {
	return i.b.IsThereChhWindows()
}

type IsHouseCloseAndALlPeriodsAtrAreEmpty struct {
	ID        string
	b         *BillingMeasure    `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
}

func NewIsHouseCloseAndALlPeriodsAtrAreEmpty(b *BillingMeasure, magnitude measures.Magnitude) *IsHouseCloseAndALlPeriodsAtrAreEmpty {
	return &IsHouseCloseAndALlPeriodsAtrAreEmpty{
		ID:        "IS_HOUSE_CLOSE_AND_ALL_ATR_PERIODS_EMPTY",
		b:         b,
		Magnitude: magnitude,
	}
}
func (i IsHouseCloseAndALlPeriodsAtrAreEmpty) Eval(_ context.Context) bool {
	return i.b.IsHouseClose() && !i.b.IsBalanceValid(i.Magnitude) && i.b.AreAllCloseAtrEmpty(i.Magnitude)
}

type IsCCCHCompleteGD struct {
	ID                         string
	b                          *BillingMeasure                             `bson:"-"`
	processedMeasureRepository process_measures.ProcessedMeasureRepository `bson:"-"`
	Magnitude                  measures.Magnitude                          `bson:"-"`
}

func NewIsCCCHCompleteGD(b *BillingMeasure, magnitude measures.Magnitude, processedMeasureRepository process_measures.ProcessedMeasureRepository) *IsCCCHCompleteGD {
	return &IsCCCHCompleteGD{
		ID:                         "IS_CCCH_COMPLETE_GD",
		b:                          b,
		processedMeasureRepository: processedMeasureRepository,
		Magnitude:                  magnitude,
	}
}

func (i IsCCCHCompleteGD) Eval(ctx context.Context) bool {

	loadCurveBillingValues := i.b.BillingLoadCurve
	curveQuarterDatesAvailable := map[string]bool{}

	if i.b.RegisterType != measures.Hourly {
		return false
	}

	for _, loadCurveBilling := range loadCurveBillingValues {
		if loadCurveBilling.Origin == measures.Filled {
			curveQuarterDatesAvailable[loadCurveBilling.EndDate.Format("2006-01-02T15")] = false
		}
	}

	loadCurves, err := i.processedMeasureRepository.ListHistoryLoadCurve(ctx, process_measures.QueryHistoryLoadCurve{
		CUPS:          i.b.CUPS,
		StartDate:     i.b.InitDate.AddDate(0, 0, -1),
		EndDate:       i.b.EndDate.AddDate(0, 0, 1),
		Periods:       i.b.Periods,
		WithCriterias: true,
		Type:          measures.QuarterHour,
		Magnitude:     i.Magnitude,
	})

	if err != nil || len(loadCurves) == 0 {
		return false
	}

	countValues := 0
	countValidValues := 0
	lastDate := loadCurves[0].EndDate.Format("2006-01-02T15")

	for _, loadCurve := range loadCurves {
		date := loadCurve.EndDate.Format("2006-01-02T15")

		if _, ok := curveQuarterDatesAvailable[date]; !ok {
			continue
		}
		countValues++

		if lastDate != date && countValues == 4 && countValidValues != 4 {
			return false
		}

		if lastDate != date {
			curveQuarterDatesAvailable[lastDate] = true

			countValues = 0
			countValidValues = 0
		}

		lastDate = date

		if loadCurve.Origin != measures.Filled {
			countValidValues++
		}

	}

	if countValues == 4 && countValidValues != 4 {
		return false
	}

	return true
}
