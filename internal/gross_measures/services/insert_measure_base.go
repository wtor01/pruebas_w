package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"time"
)

var ErrDistributor = errors.New("distributor not found")
var ErrDiffDistributors = errors.New("distributors should be the same")

type InsertMeasureBase struct {
	validationClient clients.Validation
	inventoryClient  clients.Inventory
	topic            string
	generatorDate    func() time.Time
}

func (svc InsertMeasureBase) getDistributorCDOS(measures []gross_measures.GrossMeasureBase) (string, error) {

	distributorCode := ""

	for _, m := range measures {
		if distributorCode == "" {
			distributorCode = m.GetDistributorCDOS()
		}
		if distributorCode != m.GetDistributorCDOS() {
			return distributorCode, ErrDiffDistributors
		}
	}

	if distributorCode == "" {
		return distributorCode, ErrDistributor
	}

	return distributorCode, nil
}

func (svc InsertMeasureBase) getDistributorID(measuresBase []gross_measures.GrossMeasureBase) (string, error) {

	distributorID := ""

	for _, m := range measuresBase {
		if distributorID == "" {
			distributorID = m.GetDistributorID()
		}
		if distributorID != m.GetDistributorID() {
			return distributorID, ErrDiffDistributors
		}
	}

	if distributorID == "" {
		return distributorID, ErrDistributor
	}

	return distributorID, nil
}

func (svc InsertMeasureBase) SetInsertMeasureMetadata(ctx context.Context, measuresBase []gross_measures.GrossMeasureBase) error {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "SetInsertMeasureMetadata")
	defer span.End()

	distributorCode, err := svc.getDistributorCDOS(measuresBase)

	if err != nil {
		return err
	}

	distributor, err := svc.inventoryClient.GetDistributorByCdos(ctx, distributorCode)

	if err != nil {
		return err
	}

	for i := range measuresBase {
		measuresBase[i].SetDistributorID(distributor.ID)
		measuresBase[i].GenerateID()
	}

	return nil
}

func (svc InsertMeasureBase) ValidateMeasure(ctx context.Context, measuresBase []gross_measures.GrossMeasureBase) error {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "ValidateMeasure")
	defer span.End()

	distributorID, err := svc.getDistributorID(measuresBase)

	if err != nil {
		return err
	}

	validationsConfig, err := svc.validationClient.GetValidationConfigList(ctx, distributorID, string(measures.Gross))

	if err != nil {
		return err
	}

	validators := make([]validations.ValidatorI, 0)
	for _, v := range validationsConfig {
		validator, err := validations.NewValidatorFromClient(v, nil)
		if err != nil {
			continue
		}

		validators = append(validators, validator...)
	}

	for i := range measuresBase {
		measuresBase[i].SetStatus(measures.Valid)
		for _, v := range validators {
			for _, measureValidatable := range measuresBase[i].ToValidatable() {
				validatorBase := v.Validate(measureValidatable)
				if validatorBase != nil {
					measuresBase[i].SetStatusMeasure(*validatorBase)
				}
			}
		}
	}

	return nil
}
