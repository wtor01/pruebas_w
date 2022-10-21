package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"errors"
	"time"
)

type SelfConsumptionMagnitude string

const (
	AIAutoSelfConsumptionMagnitude SelfConsumptionMagnitude = "AIAuto"
	AeAutoSelfConsumptionMagnitude SelfConsumptionMagnitude = "AeAuto"
	EHGNSelfConsumptionMagnitude   SelfConsumptionMagnitude = "EHGN"
	EHEXSelfConsumptionMagnitude   SelfConsumptionMagnitude = "EHEX"
	EHAUSelfConsumptionMagnitude   SelfConsumptionMagnitude = "EHAU"
)

type BillingLoadCurve struct {
	EndDate          time.Time                 `json:"end_date" bson:"end_date"`
	Origin           measures.OriginType       `json:"origin" bson:"origin"`
	Equipment        measures.EquipmentType    `json:"equipment" bson:"equipment"`
	Period           measures.PeriodKey        `json:"Period" bson:"Period"`
	SeasonId         string                    `json:"season_id" bson:"season_id"`
	DayTypeId        string                    `json:"day_type_id" bson:"day_type_id"`
	MeasurePointType measures.MeasurePointType `json:"measure_point_type" bson:"measure_point_type"`

	// Self Consumption

	AIAuto *float64 `json:"AI_AUTO" bson:"AI_AUTO"`
	AeAuto *float64 `json:"AE_AUTO" bson:"AE_AUTO"`
	EHGN   *float64 `json:"EHGN" bson:"EHGN"`
	EHEX   *float64 `json:"EHEX" bson:"EHEX"`
	EHAU   *float64 `json:"EHAU" bson:"EHAU"`

	// AI
	AI                       float64               `json:"AI" bson:"AI"`
	EstimatedCodeAI          int                   `json:"estimated_code_ai" bson:"estimated_code_ai"`
	EstimatedMethodAI        BalanceType           `json:"estimated_method_ai" bson:"estimated_method_ai"`
	EstimatedGeneralMethodAI GeneralEstimateMethod `json:"estimated_general_method_ai" bson:"estimated_general_method_ai"`
	MeasureTypeAI            BalanceMeasureType    `json:"measure_type_ai" bson:"measure_type_ai"`
	// AE
	AE                       float64               `json:"AE" bson:"AE"`
	EstimatedCodeAE          int                   `json:"estimated_code_ae" bson:"estimated_code_ae"`
	EstimatedMethodAE        BalanceType           `json:"estimated_method_ae" bson:"estimated_method_ae"`
	EstimatedGeneralMethodAE GeneralEstimateMethod `json:"estimated_general_method_ae" bson:"estimated_general_method_ae"`
	MeasureTypeAE            BalanceMeasureType    `json:"measure_type_ae" bson:"measure_type_ae"`
	// R1
	R1                       float64               `json:"r1" bson:"r1"`
	EstimatedCodeR1          int                   `json:"estimated_code_r1" bson:"estimated_code_r1"`
	EstimatedMethodR1        BalanceType           `json:"estimated_method_r1" bson:"estimated_method_r1"`
	EstimatedGeneralMethodR1 GeneralEstimateMethod `json:"estimated_general_method_r1" bson:"estimated_general_method_r1"`
	MeasureTypeR1            BalanceMeasureType    `json:"measure_type_r1" bson:"measure_type_r1"`
	// R2
	R2                       float64               `json:"r2" bson:"R2"`
	EstimatedCodeR2          int                   `json:"estimated_code_r2" bson:"estimated_code_r2"`
	EstimatedMethodR2        BalanceType           `json:"estimated_method_r2" bson:"estimated_method_r2"`
	EstimatedGeneralMethodR2 GeneralEstimateMethod `json:"estimated_general_method_r2" bson:"estimated_general_method_r2"`
	MeasureTypeR2            BalanceMeasureType    `json:"measure_type_r2" bson:"measure_type_r2"`
	// R3
	R3                       float64               `json:"r3" bson:"R3"`
	EstimatedCodeR3          int                   `json:"estimated_code_r3" bson:"estimated_code_r3"`
	EstimatedMethodR3        BalanceType           `json:"estimated_method_r3" bson:"estimated_method_r3"`
	EstimatedGeneralMethodR3 GeneralEstimateMethod `json:"estimated_general_method_r3" bson:"estimated_general_method_r3"`
	MeasureTypeR3            BalanceMeasureType    `json:"measure_type_r3" bson:"measure_type_r3"`
	// R4
	R4                       float64               `json:"r4" bson:"R4"`
	EstimatedCodeR4          int                   `json:"estimated_code_r4" bson:"estimated_code_r4"`
	EstimatedMethodR4        BalanceType           `json:"estimated_method_r4" bson:"estimated_method_r4"`
	EstimatedGeneralMethodR4 GeneralEstimateMethod `json:"estimated_general_method_r4" bson:"estimated_general_method_r4"`
	MeasureTypeR4            BalanceMeasureType    `json:"measure_type_r4" bson:"measure_type_r4"`
}

func (c *BillingLoadCurve) setEstimatedCode(magnitude measures.Magnitude, estimateCode int) {
	switch magnitude {
	case measures.AI:
		c.EstimatedCodeAI = estimateCode
	case measures.AE:
		c.EstimatedCodeAE = estimateCode
	case measures.R1:
		c.EstimatedCodeR1 = estimateCode
	case measures.R2:
		c.EstimatedCodeR2 = estimateCode
	case measures.R3:
		c.EstimatedCodeR3 = estimateCode
	case measures.R4:
		c.EstimatedCodeR4 = estimateCode
	}
}

func (c *BillingLoadCurve) setEstimatedMethod(magnitude measures.Magnitude, estimatedMethod BalanceType) {
	switch magnitude {
	case measures.AI:
		c.EstimatedMethodAI = estimatedMethod
	case measures.AE:
		c.EstimatedMethodAE = estimatedMethod
	case measures.R1:
		c.EstimatedMethodR1 = estimatedMethod
	case measures.R2:
		c.EstimatedMethodR2 = estimatedMethod
	case measures.R3:
		c.EstimatedMethodR3 = estimatedMethod
	case measures.R4:
		c.EstimatedMethodR4 = estimatedMethod
	}
}

func (c *BillingLoadCurve) setMeasureType(magnitude measures.Magnitude, measureType BalanceMeasureType) {
	switch magnitude {
	case measures.AI:
		c.MeasureTypeAI = measureType
	case measures.AE:
		c.MeasureTypeAE = measureType
	case measures.R1:
		c.MeasureTypeR1 = measureType
	case measures.R2:
		c.MeasureTypeR2 = measureType
	case measures.R3:
		c.MeasureTypeR3 = measureType
	case measures.R4:
		c.MeasureTypeR4 = measureType
	}
}

func (c *BillingLoadCurve) GetMagnitude(magnitude measures.Magnitude) float64 {
	switch magnitude {
	case measures.AI:
		return c.AI
	case measures.AE:
		return c.AE
	case measures.R1:
		return c.R1
	case measures.R2:
		return c.R2
	case measures.R3:
		return c.R3
	case measures.R4:
		return c.R4
	default:
		return .0
	}
}

func (c *BillingLoadCurve) SumMagnitude(magnitude measures.Magnitude, value float64) {
	switch magnitude {
	case measures.AI:
		c.AI += value
	case measures.AE:
		c.AE += value
	case measures.R1:
		c.R1 += value
	case measures.R2:
		c.R2 += value
	case measures.R3:
		c.R3 += value
	case measures.R4:
		c.R4 += value
	}
}

func (c *BillingLoadCurve) SetMagnitude(magnitude measures.Magnitude, value float64) {
	switch magnitude {
	case measures.AI:
		c.AI = value
	case measures.AE:
		c.AE = value
	case measures.R1:
		c.R1 = value
	case measures.R2:
		c.R2 = value
	case measures.R3:
		c.R3 = value
	case measures.R4:
		c.R4 = value
	}
}

func (c *BillingLoadCurve) SumSelfConsumptionMagnitude(magnitude SelfConsumptionMagnitude, value *float64) {
	if value == nil {
		return
	}

	generateNewValue := func(oldValuePtr *float64) float64 {
		oldValue := .0
		if oldValuePtr != nil {
			oldValue = *oldValuePtr
		}
		return oldValue + *value
	}

	switch magnitude {
	case AIAutoSelfConsumptionMagnitude:
		newValue := generateNewValue(c.AIAuto)
		c.AIAuto = &newValue
	case AeAutoSelfConsumptionMagnitude:
		newValue := generateNewValue(c.AeAuto)
		c.AeAuto = &newValue
	case EHAUSelfConsumptionMagnitude:
		newValue := generateNewValue(c.EHAU)
		c.EHAU = &newValue
	case EHGNSelfConsumptionMagnitude:
		newValue := generateNewValue(c.EHGN)
		c.EHGN = &newValue
	case EHEXSelfConsumptionMagnitude:
		newValue := generateNewValue(c.EHEX)
		c.EHEX = &newValue
	}
}

func (c *BillingLoadCurve) SetSelfConsumptionMagnitude(magnitude SelfConsumptionMagnitude, value *float64) {
	if value == nil {
		return
	}

	switch magnitude {
	case AIAutoSelfConsumptionMagnitude:
		c.AIAuto = value
	case AeAutoSelfConsumptionMagnitude:
		c.AeAuto = value
	case EHAUSelfConsumptionMagnitude:
		c.EHAU = value
	case EHGNSelfConsumptionMagnitude:
		c.EHGN = value
	case EHEXSelfConsumptionMagnitude:
		c.EHEX = value
	}
}

func (c BillingLoadCurve) getOriginType() measures.OriginType {
	return c.Origin
}
func (c BillingLoadCurve) getPeriod() measures.PeriodKey {
	return c.Period
}

func (c BillingLoadCurve) getEstimateCode(magnitude measures.Magnitude) (int, error) {
	switch magnitude {
	case measures.AI:
		return c.EstimatedCodeAI, nil
	case measures.AE:
		return c.EstimatedCodeAE, nil
	case measures.R1:
		return c.EstimatedCodeR1, nil
	case measures.R2:
		return c.EstimatedCodeR2, nil
	case measures.R3:
		return c.EstimatedCodeR3, nil
	case measures.R4:
		return c.EstimatedCodeR4, nil
	default:
		return 0, errors.New("unexpected magnitude")
	}
}

func (c BillingLoadCurve) getGeneralEstimatedMethod(magnitude measures.Magnitude) (GeneralEstimateMethod, error) {
	switch magnitude {
	case measures.AI:
		return c.EstimatedGeneralMethodAI, nil
	case measures.AE:
		return c.EstimatedGeneralMethodAE, nil
	case measures.R1:
		return c.EstimatedGeneralMethodR1, nil
	case measures.R2:
		return c.EstimatedGeneralMethodR2, nil
	case measures.R3:
		return c.EstimatedGeneralMethodR3, nil
	case measures.R4:
		return c.EstimatedGeneralMethodR4, nil
	default:
		return GeneralOutlined, errors.New("unexpected magnitude")
	}
}

func (c *BillingLoadCurve) setGeneralEstimatedMethod(magnitude measures.Magnitude, method GeneralEstimateMethod) {
	switch magnitude {
	case measures.AI:
		c.EstimatedGeneralMethodAI = method
	case measures.AE:
		c.EstimatedGeneralMethodAE = method
	case measures.R1:
		c.EstimatedGeneralMethodR1 = method
	case measures.R2:
		c.EstimatedGeneralMethodR2 = method
	case measures.R3:
		c.EstimatedGeneralMethodR3 = method
	case measures.R4:
		c.EstimatedGeneralMethodR4 = method
	}
}

func (c BillingLoadCurve) getMagnitude(magnitude measures.Magnitude) (float64, error) {
	switch magnitude {
	case measures.AI:
		return c.AI, nil
	case measures.AE:
		return c.AE, nil
	case measures.R1:
		return c.R1, nil
	case measures.R2:
		return c.R2, nil
	case measures.R3:
		return c.R3, nil
	case measures.R4:
		return c.R4, nil
	default:
		return .0, errors.New("unexpected magnitude")
	}
}
