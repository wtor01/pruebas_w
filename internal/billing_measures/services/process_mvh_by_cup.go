package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

type ProcessMvhByCup struct {
	repository                 billing_measures.BillingMeasureRepository
	repositoryProfiles         billing_measures.ConsumProfileRepository
	processedMeasureRepository process_measures.ProcessedMeasureRepository
	inventoryClient            clients.Inventory
	calendarPeriodRepository   measures.CalendarPeriodRepository
	Location                   *time.Location
	masterTablesClient         clients.MasterTables
	tracer                     trace.Tracer
	publisher                  event.PublisherCreator
	topic                      string
}

func NewProcessMvhByCup(
	repository billing_measures.BillingMeasureRepository,
	processedMeasureRepository process_measures.ProcessedMeasureRepository,
	inventoryClient clients.Inventory,
	repositoryProfiles billing_measures.ConsumProfileRepository,
	calendarPeriodRepository measures.CalendarPeriodRepository,
	Location *time.Location,
	masterTablesClient clients.MasterTables,
	publisher event.PublisherCreator,
	topic string,
) *ProcessMvhByCup {
	return &ProcessMvhByCup{
		repository:                 repository,
		processedMeasureRepository: processedMeasureRepository,
		inventoryClient:            inventoryClient,
		repositoryProfiles:         repositoryProfiles,
		calendarPeriodRepository:   calendarPeriodRepository,
		Location:                   Location,
		masterTablesClient:         masterTablesClient,
		tracer:                     telemetry.GetTracer(),
		publisher:                  publisher,
		topic:                      topic,
	}
}

func (svc ProcessMvhByCup) Handler(ctx context.Context, dto measures.ProcessMeasurePayload) error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - Handler")
	defer span.End()

	tariff, err := svc.masterTablesClient.GetTariff(ctx, clients.GetTariffDto{
		ID: dto.MeterConfig.TariffID(),
	})

	if err != nil {
		return err
	}

	calendarPeriodDay, err := svc.calendarPeriodRepository.GetCalendarPeriod(ctx, measures.SearchCalendarPeriod{
		Day:          dto.Date,
		GeographicID: tariff.GeographicId,
		CalendarCode: tariff.CalendarId,
		Location:     svc.Location,
	})

	if err != nil {
		return err
	}

	periods := calendarPeriodDay.GetAllPeriods()

	magnitudes := dto.MeterConfig.GetMagnitudesActive()

	lastBillingMeasure, err := svc.repository.Last(ctx, billing_measures.QueryLast{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date,
	})

	billingMeasure := billing_measures.NewBillingMeasure(
		dto.MeterConfig.Cups(),
		lastBillingMeasure.EndDate,
		dto.Date.AddDate(0, 0, 1),
		lastBillingMeasure.DistributorCode,
		lastBillingMeasure.DistributorID,
		periods,
		magnitudes,
		dto.MeterConfig.Type,
	)

	// TODO: cambair el NewBillingMeasure para crearse a artir del dto.MeterConfig
	billingMeasure.Inaccessible = dto.MeterConfig.Inaccessible
	billingMeasure.PointType = dto.MeterConfig.PointType()
	billingMeasure.RegisterType = dto.MeterConfig.CurveType
	billingMeasure.Technology = dto.MeterConfig.Meter.Technology

	//SUPERVISION
	if err != nil || lastBillingMeasure.EndDate.IsZero() {
		errSave := svc.setSupervision(ctx, billing_measures.NoLastBillingMeasure, billingMeasure)
		if errSave != nil {
			return errSave
		}
		return err
	}

	previous, _ := svc.repository.GetPrevious(ctx, billing_measures.GetPrevious{
		CUPS:     billingMeasure.CUPS,
		InitDate: billingMeasure.InitDate,
		EndDate:  billingMeasure.EndDate,
	})

	billingMeasure.CalculateVersionByPreviousBillingMeasure(previous)

	billingMeasure.SetContractInfo(dto.MeterConfig)
	billingMeasure.SetCoefficient(tariff)

	billingMeasure.Status = billing_measures.Calculating

	err = svc.repository.Save(ctx, billingMeasure)

	if err != nil {
		return err
	}

	err = svc.setBillingLoadCurve(ctx, dto, &lastBillingMeasure, &billingMeasure)

	if err != nil {
		return err
	}

	err = svc.setReadingsClosure(ctx, dto, &lastBillingMeasure, &billingMeasure)

	if err != nil {
		//SUPERVISION
		err = svc.setSupervision(ctx, billing_measures.NotValidReadingClosure, billingMeasure)
		return err
	}

	billingMeasure.CalcAtrBalance()

	billingMeasure.CalcAtrVsCurve()

	if billingMeasure.ShouldExecuteGraph() {
		processErr := svc.processGraph(ctx, dto, &billingMeasure)
		billingMeasure.Calculated()
		if len(processErr) > 0 {
			billingMeasure.Supervision()
			errDesc := strings.Join(utils.MapSlice(processErr, func(item error) string {
				return item.Error()
			}), ",")
			billingMeasure.SetDescriptionStatus(billing_measures.DescriptionStatus(errDesc))
		}
		billingMeasure.AfterExecuteGraph()
	}

	billingMeasure.BeforeSave(dto.MeterConfig)

	err = svc.repository.Save(ctx, billingMeasure)

	if err != nil {
		return err
	}

	err = event.PublishEvent(ctx, svc.topic, svc.publisher, billing_measures.NewOnSaveBillingMeasureEvent(billing_measures.OnSaveBillingMeasurePayload{
		InitDate:         billingMeasure.InitDate,
		EndDate:          billingMeasure.EndDate,
		CUPS:             billingMeasure.CUPS,
		BillingMeasureId: billingMeasure.Id,
	}))

	return err
}

func (svc ProcessMvhByCup) setSupervision(ctx context.Context, description billing_measures.DescriptionStatus, billingMeasure billing_measures.BillingMeasure) error {
	billingMeasure.Supervision()
	billingMeasure.SetDescriptionStatus(description)
	err := svc.repository.Save(ctx, billingMeasure)
	return err
}

func (svc ProcessMvhByCup) processTlgGraph(
	ctx context.Context,
	billingMeasure *billing_measures.BillingMeasure,
) []error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - processTlgGraph")
	defer span.End()
	err := make([]error, 0, len(billingMeasure.Periods)*len(billingMeasure.Magnitudes)*2)
	for _, period := range billingMeasure.Periods {
		for _, magnitude := range billingMeasure.Magnitudes {
			graph := billing_measures.GenerateTree(billingMeasure, period, magnitude, svc.repository, svc.repositoryProfiles)
			singleErr := graph.Execute(ctx)
			if singleErr != nil {
				err = append(err, singleErr)
			}
			key := fmt.Sprintf("%s_%s", period, magnitude)
			billingMeasure.GraphHistory[key] = graph
		}
	}
	return err
}

func (svc ProcessMvhByCup) processDcNoTlgGraph(
	ctx context.Context,
	billingMeasure *billing_measures.BillingMeasure,
) []error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - processDcNoTlgGraph")
	defer span.End()
	err := make([]error, 0, len(billingMeasure.Magnitudes)*2)
	for _, magnitude := range billingMeasure.Magnitudes {
		graph := billing_measures.GenerateDcTreeNoTlg(billingMeasure, svc.processedMeasureRepository, magnitude, svc.repository, svc.repositoryProfiles)
		singleErr := graph.Execute(ctx)
		if singleErr != nil {
			err = append(err, singleErr)
		}
		key := fmt.Sprintf("%s", magnitude)
		billingMeasure.GraphHistory[key] = graph
	}
	return err
}

func (svc ProcessMvhByCup) processGdNoTlgGraph(
	ctx context.Context,
	billingMeasure *billing_measures.BillingMeasure,
) []error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - processGdNoTlgGraph")
	defer span.End()
	err := make([]error, 0, len(billingMeasure.Magnitudes)*2)
	for _, magnitude := range billingMeasure.Magnitudes {
		graph := billing_measures.GenerateGDTreeNoTlg(billingMeasure, magnitude, svc.repositoryProfiles, svc.repository, svc.processedMeasureRepository)
		singleErr := graph.Execute(ctx)
		if singleErr != nil {
			err = append(err, singleErr)
		}
		key := fmt.Sprintf("%s", magnitude)
		billingMeasure.GraphHistory[key] = graph
	}
	return err
}

func (svc ProcessMvhByCup) processGraph(
	ctx context.Context,
	dto measures.ProcessMeasurePayload,
	billingMeasure *billing_measures.BillingMeasure,
) []error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - processGraph")
	defer span.End()

	span.SetAttributes(attribute.String("service_type", string(dto.MeterConfig.ServiceType())))
	span.SetAttributes(attribute.String("meter_config.type", string(dto.MeterConfig.Type)))

	if dto.MeterConfig.ServiceType() == measures.DcServiceType {
		if dto.MeterConfig.Type == measures.TLG {
			return svc.processTlgGraph(ctx, billingMeasure)

		}
		return svc.processDcNoTlgGraph(ctx, billingMeasure)
	}
	if dto.MeterConfig.ServiceType() == measures.GdServiceType {
		if dto.MeterConfig.Type == measures.TLG {
			return svc.processTlgGraph(ctx, billingMeasure)

		}
		return svc.processGdNoTlgGraph(ctx, billingMeasure)
	}
	return []error{errors.New("INVALID SERVICE TYPE")}
}

func (svc ProcessMvhByCup) getDailyReadingClosureByCups(ctx context.Context, dto measures.ProcessMeasurePayload) (measures.DailyReadingClosure, error) {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - getDailyReadingClosureByCups")
	defer span.End()

	monthly, err := svc.processedMeasureRepository.GetMonthlyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date.AddDate(0, 0, 1),
	})

	if err == nil {
		return monthly.ToDailyReadingClosure(), nil
	}

	daily, err := svc.processedMeasureRepository.GetProcessedDailyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date.AddDate(0, 0, 1),
	})

	if err != nil {
		return measures.DailyReadingClosure{}, err
	}

	return daily.ToDailyReadingClosure(), nil
}

func (svc ProcessMvhByCup) setReadingsClosure(
	ctx context.Context,
	dto measures.ProcessMeasurePayload,
	lastBillingMeasure *billing_measures.BillingMeasure,
	billingMeasure *billing_measures.BillingMeasure,
) error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - setReadingsClosure")
	defer span.End()

	billingMeasure.SetPreviousReadingClosure(lastBillingMeasure.ActualReadingClosure)

	readingDate, err := svc.getDailyReadingClosureByCups(ctx, dto)

	if err == nil {
		if !readingDate.InitDate.IsZero() && readingDate.InitDate.Before(lastBillingMeasure.EndDate) {
			return errors.New("not valid reading closure")
		}
		billingMeasure.SetActualReadingClosure(readingDate)
	}

	billingMeasure.SetExecutionBalanceOrigin(readingDate)
	return nil
}

func (svc ProcessMvhByCup) setBillingLoadCurve(
	ctx context.Context,
	dto measures.ProcessMeasurePayload,
	lastBillingMeasure *billing_measures.BillingMeasure,
	b *billing_measures.BillingMeasure,
) error {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - setBillingLoadCurve")
	defer span.End()

	processedLoadCurve, err := svc.processedMeasureRepository.ProcessedLoadCurveByCups(ctx, process_measures.QueryProcessedLoadCurveByCups{
		CUPS:      dto.MeterConfig.Cups(),
		StartDate: lastBillingMeasure.EndDate,
		EndDate:   dto.Date.AddDate(0, 0, 1),
		CurveType: measures.HourlyMeasureCurveReadingType,
		Status:    measures.Valid,
	})

	if err != nil {
		return err
	}

	processedLoadCurve, err = svc.fillEmptyCurves(
		ctx,
		dto.MeterConfig.Cups(),
		lastBillingMeasure.EndDate,
		dto.Date.AddDate(0, 0, 1),
		processedLoadCurve,
	)

	if err != nil {
		return err
	}

	b.SetBillingLoadCurve(utils.MapSlice(processedLoadCurve, func(item process_measures.ProcessedLoadCurve) billing_measures.BillingLoadCurve {
		return billing_measures.BillingLoadCurve{
			EndDate:          item.EndDate,
			Origin:           item.Origin,
			AI:               item.AI,
			AE:               item.AE,
			R1:               item.R1,
			R2:               item.R2,
			R3:               item.R3,
			R4:               item.R4,
			Period:           item.Period,
			MeasurePointType: item.MeasurePointType,
		}
	}))

	return nil
}

func (svc ProcessMvhByCup) fillEmptyCurves(
	ctx context.Context,
	cups string,
	startDate, endDate time.Time,
	processedLoadCurve []process_measures.ProcessedLoadCurve,
) ([]process_measures.ProcessedLoadCurve, error) {
	ctx, span := svc.tracer.Start(ctx, "ProcessMvhByCup - fillEmptyCurves")
	defer span.End()

	curvesFilled := make(map[string]struct {
		curve        *process_measures.ProcessedLoadCurve
		curveOrigin  *process_measures.ProcessedLoadCurve
		quarterDates int
	})

	generateKeyDate := func(date time.Time) string {
		return date.Format("2006-01-02 15")
	}

	for i, curve := range processedLoadCurve {
		if curve.Origin != measures.Filled {
			continue
		}

		curvesFilled[generateKeyDate(curve.EndDate)] = struct {
			curve        *process_measures.ProcessedLoadCurve
			curveOrigin  *process_measures.ProcessedLoadCurve
			quarterDates int
		}{curveOrigin: &processedLoadCurve[i], curve: &curve, quarterDates: 0}
	}

	if len(curvesFilled) == 0 {
		return processedLoadCurve, nil
	}

	processedLoadCurveQuarter, err := svc.processedMeasureRepository.ProcessedLoadCurveByCups(ctx, process_measures.QueryProcessedLoadCurveByCups{
		CUPS:      cups,
		StartDate: startDate,
		EndDate:   endDate,
		CurveType: measures.QuarterMeasureCurveReadingType,
		Status:    measures.Valid,
	})

	if err != nil {
		return []process_measures.ProcessedLoadCurve{}, err
	}

	for _, curve := range processedLoadCurveQuarter {
		if curve.Origin == measures.Filled {
			continue
		}
		date := curve.EndDate

		if date.Hour() != 0 {
			// añadimos 45 minutos para que se cumpla esto:
			// las medidas cuartohorarias -> 11:15 11:30 11:45 12:00, son de las hora de las 12:00
			date = date.Add(time.Minute * 45)
		}

		keyDate := generateKeyDate(date)
		curveInfo, ok := curvesFilled[keyDate]

		if !ok {
			continue
		}

		curveInfo.curve.AI += curve.AI
		curveInfo.curve.AE += curve.AE
		curveInfo.curve.R1 += curve.R1
		curveInfo.curve.R2 += curve.R2
		curveInfo.curve.R3 += curve.R3
		curveInfo.curve.R4 += curve.R4
		curveInfo.quarterDates += 1

		if curveInfo.quarterDates == 4 {
			curveInfo.curveOrigin.AI = curveInfo.curve.AI
			curveInfo.curveOrigin.AE = curveInfo.curve.AE
			curveInfo.curveOrigin.R1 = curveInfo.curve.R1
			curveInfo.curveOrigin.R2 = curveInfo.curve.R2
			curveInfo.curveOrigin.R3 = curveInfo.curve.R3
			curveInfo.curveOrigin.R4 = curveInfo.curve.R4
			curveInfo.curveOrigin.Origin = measures.CalculatedWithQuarter
		}

		curvesFilled[keyDate] = curveInfo

	}

	return processedLoadCurve, nil
}
