package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

type ProcessMeasureBase struct {
	generatorDate            func() time.Time
	calendarPeriodRepository measures.CalendarPeriodRepository
	Location                 *time.Location
	validationClient         clients.Validation
	masterTablesClient       clients.MasterTables
}

func NewProcessMeasureBase(
	calendarPeriodRepository measures.CalendarPeriodRepository,
	Location *time.Location,
	validationClient clients.Validation,
	masterTablesClient clients.MasterTables,

) ProcessMeasureBase {
	return ProcessMeasureBase{
		generatorDate:            time.Now,
		calendarPeriodRepository: calendarPeriodRepository,
		Location:                 Location,
		validationClient:         validationClient,
		masterTablesClient:       masterTablesClient,
	}
}

func (svc ProcessMeasureBase) GetCalendarPeriod(
	ctx context.Context,
	dto measures.ProcessMeasurePayload,
	tariff clients.Tariffs,
) (measures.CalendarPeriod, error) {

	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetCalendarPeriod")
	defer span.End()

	calendar, err := svc.calendarPeriodRepository.GetCalendarPeriod(ctx, measures.SearchCalendarPeriod{
		Day:          dto.Date,
		GeographicID: tariff.GeographicId,
		CalendarCode: tariff.CalendarId,
		Location:     svc.Location,
	})

	return calendar, err
}

func (svc ProcessMeasureBase) GetTariff(ctx context.Context, dto measures.ProcessMeasurePayload) (clients.Tariffs, error) {

	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetTariff")
	defer span.End()

	tariff, err := svc.masterTablesClient.GetTariff(ctx, clients.GetTariffDto{
		ID: dto.MeterConfig.TariffID(),
	})

	return tariff, err
}

func (svc ProcessMeasureBase) GetValidations(ctx context.Context, dto measures.ProcessMeasurePayload, validationRepository validations.ValidationMongoRepository) ([]validations.ValidatorI, error) {

	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetValidations")
	defer span.End()

	span.SetAttributes(attribute.String("validationType", string(measures.Process)))

	validationsConfig, err := svc.validationClient.GetValidationConfigList(ctx, dto.MeterConfig.DistributorID, string(measures.Process))
	validators := make([]validations.ValidatorI, 0)

	if err != nil {
		return validators, err
	}

	for _, v := range validationsConfig {
		validator, err := validations.NewValidatorFromClient(v, validationRepository)
		if err != nil {
			continue
		}
		validators = append(validators, validator...)
	}

	return validators, err
}

func (svc ProcessMeasureBase) ValidateMeasures(processMeasures process_measures.ToMeasuresBase, validators []validations.ValidatorI) {
	measuresBase := processMeasures.MeasuresBase()

	for i := range measuresBase {
		svc.ValidateMeasure(measuresBase[i], validators)
	}
}

func (svc ProcessMeasureBase) ValidateMeasure(processMeasures process_measures.ProcessMeasureBase, validators []validations.ValidatorI) process_measures.ProcessMeasureBase {
	for _, v := range validators {
		validatorBase := v.Validate(processMeasures.ToValidatable())
		if validatorBase != nil {
			processMeasures.SetStatusMeasure(*validatorBase)
		}
	}

	return processMeasures
}
