package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

type BillingBalance struct {
	EndDate time.Time             `json:"end_date" bson:"end_date,omitempty"`
	Origin  measures.OriginType   `json:"origin "bson:"origin,omitempty"`
	P0      *BillingBalancePeriod `json:"p0" bson:"p0,omitempty"`
	P1      *BillingBalancePeriod `json:"p1" bson:"p1,omitempty"`
	P2      *BillingBalancePeriod `json:"p2" bson:"p2,omitempty"`
	P3      *BillingBalancePeriod `json:"p3" bson:"p3,omitempty"`
	P4      *BillingBalancePeriod `json:"p4" bson:"p4,omitempty"`
	P5      *BillingBalancePeriod `json:"p5" bson:"p5,omitempty"`
	P6      *BillingBalancePeriod `json:"p6" bson:"p6,omitempty"`
}

func (balance BillingBalance) getBillingBalancePeriod(p measures.PeriodKey) BillingBalancePeriod {

	balancePeriod := &BillingBalancePeriod{}

	switch p {
	case measures.P0:
		balancePeriod = balance.P0
	case measures.P1:
		balancePeriod = balance.P1
	case measures.P2:
		balancePeriod = balance.P2
	case measures.P3:
		balancePeriod = balance.P3
	case measures.P4:
		balancePeriod = balance.P4
	case measures.P5:
		balancePeriod = balance.P5
	case measures.P6:
		balancePeriod = balance.P6
	}

	return *balancePeriod
}

type BillingBalancePeriod struct {

	// AI
	AI                   float64               `json:"AI" bson:"AI"`
	BalanceTypeAI        BalanceType           `json:"balance_type_ai" bson:"balance_type_ai"`
	BalanceGeneralTypeAI GeneralEstimateMethod `json:"balance_general_type_ai" bson:"balance_general_type_ai"`
	BalanceMeasureTypeAI BalanceMeasureType    `json:"balance_measure_type_ai" bson:"balance_measure_type_ai"`
	BalanceOriginAI      BalanceOriginType     `json:"balance_origin_ai" bson:"balance_origin_ai"`
	BalanceValidationAI  measures.Status       `json:"balance_validation_ai" bson:"balance_validation_ai"`
	EstimatedCodeAI      int                   `json:"estimated_code_ai" bson:"estimated_code_ai"`

	// AE
	AE                   float64               `json:"AE" bson:"AE"`
	BalanceTypeAE        BalanceType           `json:"balance_type_ae" bson:"balance_type_ae"`
	BalanceGeneralTypeAE GeneralEstimateMethod `json:"balance_general_type_ae" bson:"balance_general_type_ae"`
	BalanceMeasureTypeAE BalanceMeasureType    `json:"balance_measure_type_ae" bson:"balance_measure_type_ae"`
	BalanceOriginAE      BalanceOriginType     `json:"balance_origin_ae" bson:"balance_origin_ae"`
	BalanceValidationAE  measures.Status       `json:"balance_validation_ae" bson:"balance_validation_ae"`
	EstimatedCodeAE      int                   `json:"estimated_code_ae" bson:"estimated_code_ae"`

	// R1
	R1                   float64               `json:"r1" bson:"R1"`
	BalanceTypeR1        BalanceType           `json:"balance_type_r1" bson:"balance_type_r1"`
	BalanceGeneralTypeR1 GeneralEstimateMethod `json:"balance_general_type_r1" bson:"balance_general_type_r1"`
	BalanceMeasureTypeR1 BalanceMeasureType    `json:"balance_measure_type_r1" bson:"balance_measure_type_r1"`
	BalanceOriginR1      BalanceOriginType     `json:"balance_origin_r1" bson:"balance_origin_r1"`
	BalanceValidationR1  measures.Status       `json:"balance_validation_r1" bson:"balance_validation_r1"`
	EstimatedCodeR1      int                   `json:"estimated_code_r1" bson:"estimated_code_r1"`

	// R2
	R2                   float64               `json:"r2" bson:"R2"`
	BalanceTypeR2        BalanceType           `json:"balance_type_r2" bson:"balance_type_r2"`
	BalanceGeneralTypeR2 GeneralEstimateMethod `json:"balance_general_type_r2" bson:"balance_general_type_r2"`
	BalanceMeasureTypeR2 BalanceMeasureType    `json:"balance_measure_type_r2" bson:"balance_measure_type_r2"`
	BalanceOriginR2      BalanceOriginType     `json:"balance_origin_r2" bson:"balance_origin_r2"`
	BalanceValidationR2  measures.Status       `json:"balance_validation_r2" bson:"balance_validation_r2"`
	EstimatedCodeR2      int                   `json:"estimated_code_r2" bson:"estimated_code_r2"`

	// R3
	R3                   float64               `json:"r3" bson:"R3"`
	BalanceTypeR3        BalanceType           `json:"balance_type_r3" bson:"balance_type_r3"`
	BalanceGeneralTypeR3 GeneralEstimateMethod `json:"balance_general_type_r3" bson:"balance_general_type_r3"`
	BalanceMeasureTypeR3 BalanceMeasureType    `json:"balance_measure_type_r3" bson:"balance_measure_type_r3"`
	BalanceOriginR3      BalanceOriginType     `json:"balance_origin_r3" bson:"balance_origin_r3"`
	BalanceValidationR3  measures.Status       `json:"balance_validation_r3" bson:"balance_validation_r3"`
	EstimatedCodeR3      int                   `json:"estimated_code_r3" bson:"estimated_code_r3"`

	// R4
	R4                   float64               `json:"r4" bson:"R4"`
	BalanceTypeR4        BalanceType           `json:"balance_type_r4" bson:"balance_type_r4"`
	BalanceGeneralTypeR4 GeneralEstimateMethod `json:"balance_general_type_r4" bson:"balance_general_type_r4"`
	BalanceMeasureTypeR4 BalanceMeasureType    `json:"balance_measure_type_r4" bson:"balance_measure_type_r4"`
	BalanceOriginR4      BalanceOriginType     `json:"balance_origin_r4" bson:"balance_origin_r4"`
	BalanceValidationR4  measures.Status       `json:"balance_validation_r4" bson:"balance_validation_r4"`
	EstimatedCodeR4      int                   `json:"estimated_code_r4" bson:"estimated_code_r4"`
}

func (balancePeriod *BillingBalancePeriod) setEstimateCode(magnitude measures.Magnitude, estimateCode int) {
	switch magnitude {
	case measures.AI:
		balancePeriod.EstimatedCodeAI = estimateCode
	case measures.AE:
		balancePeriod.EstimatedCodeAE = estimateCode
	case measures.R1:
		balancePeriod.EstimatedCodeR1 = estimateCode
	case measures.R2:
		balancePeriod.EstimatedCodeR2 = estimateCode
	case measures.R3:
		balancePeriod.EstimatedCodeR3 = estimateCode
	case measures.R4:
		balancePeriod.EstimatedCodeR4 = estimateCode
	}
}

func (balancePeriod *BillingBalancePeriod) getEstimateCode(magnitude measures.Magnitude) int {
	switch magnitude {
	case measures.AI:
		return balancePeriod.EstimatedCodeAI
	case measures.AE:
		return balancePeriod.EstimatedCodeAE
	case measures.R1:
		return balancePeriod.EstimatedCodeR1
	case measures.R2:
		return balancePeriod.EstimatedCodeR2
	case measures.R3:
		return balancePeriod.EstimatedCodeR3
	case measures.R4:
		return balancePeriod.EstimatedCodeR4
	default:
		return 0
	}
}

func (balancePeriod *BillingBalancePeriod) setBalanceType(magnitude measures.Magnitude, balanceType BalanceType) {
	switch magnitude {
	case measures.AI:
		balancePeriod.BalanceTypeAI = balanceType
	case measures.AE:
		balancePeriod.BalanceTypeAE = balanceType
	case measures.R1:
		balancePeriod.BalanceTypeR1 = balanceType
	case measures.R2:
		balancePeriod.BalanceTypeR2 = balanceType
	case measures.R3:
		balancePeriod.BalanceTypeR3 = balanceType
	case measures.R4:
		balancePeriod.BalanceTypeR4 = balanceType
	}
}

func (balancePeriod *BillingBalancePeriod) setBalanceMeasureType(magnitude measures.Magnitude, measureType BalanceMeasureType) {
	switch magnitude {
	case measures.AI:
		balancePeriod.BalanceMeasureTypeAI = measureType
	case measures.AE:
		balancePeriod.BalanceMeasureTypeAE = measureType
	case measures.R1:
		balancePeriod.BalanceMeasureTypeR1 = measureType
	case measures.R2:
		balancePeriod.BalanceMeasureTypeR2 = measureType
	case measures.R3:
		balancePeriod.BalanceMeasureTypeR3 = measureType
	case measures.R4:
		balancePeriod.BalanceMeasureTypeR4 = measureType
	}
}

func (balancePeriod *BillingBalancePeriod) getBalanceOrigin(magnitude measures.Magnitude) BalanceOriginType {
	switch magnitude {
	case measures.AI:
		return balancePeriod.BalanceOriginAI
	case measures.AE:
		return balancePeriod.BalanceOriginAE
	case measures.R1:
		return balancePeriod.BalanceOriginR1
	case measures.R2:
		return balancePeriod.BalanceOriginR2
	case measures.R3:
		return balancePeriod.BalanceOriginR3
	case measures.R4:
		return balancePeriod.BalanceOriginR4
	default:
		return ""
	}
}

func (balancePeriod *BillingBalancePeriod) setBalanceOrigin(magnitude measures.Magnitude, val BalanceOriginType) {
	switch magnitude {
	case measures.AI:
		balancePeriod.BalanceOriginAI = val
	case measures.AE:
		balancePeriod.BalanceOriginAE = val
	case measures.R1:
		balancePeriod.BalanceOriginR1 = val
	case measures.R2:
		balancePeriod.BalanceOriginR2 = val
	case measures.R3:
		balancePeriod.BalanceOriginR3 = val
	case measures.R4:
		balancePeriod.BalanceOriginR4 = val
	}
}

func (balancePeriod *BillingBalancePeriod) getBalanceValidation(magnitude measures.Magnitude) measures.Status {
	switch magnitude {
	case measures.AI:
		return balancePeriod.BalanceValidationAI
	case measures.AE:
		return balancePeriod.BalanceValidationAE
	case measures.R1:
		return balancePeriod.BalanceValidationR1
	case measures.R2:
		return balancePeriod.BalanceValidationR2
	case measures.R3:
		return balancePeriod.BalanceValidationR3
	case measures.R4:
		return balancePeriod.BalanceValidationR4
	default:
		return measures.Invalid
	}
}

func (balancePeriod *BillingBalancePeriod) GetMagnitude(magnitude measures.Magnitude) float64 {
	switch magnitude {
	case measures.AI:
		return balancePeriod.AI
	case measures.AE:
		return balancePeriod.AE
	case measures.R1:
		return balancePeriod.R1
	case measures.R2:
		return balancePeriod.R2
	case measures.R3:
		return balancePeriod.R3
	case measures.R4:
		return balancePeriod.R4
	default:
		return .0
	}
}

func (balancePeriod *BillingBalancePeriod) GetStatus(magnitude measures.Magnitude) measures.Status {
	switch magnitude {
	case measures.AI:
		return balancePeriod.BalanceValidationAI
	case measures.AE:
		return balancePeriod.BalanceValidationAE
	case measures.R1:
		return balancePeriod.BalanceValidationR1
	case measures.R2:
		return balancePeriod.BalanceValidationR2
	case measures.R3:
		return balancePeriod.BalanceValidationR3
	case measures.R4:
		return balancePeriod.BalanceValidationR4
	default:
		return ""
	}
}

func (balancePeriod *BillingBalancePeriod) SetStatus(magnitude measures.Magnitude, status measures.Status) {
	switch magnitude {
	case measures.AI:
		balancePeriod.BalanceValidationAI = status
	case measures.AE:
		balancePeriod.BalanceValidationAE = status
	case measures.R1:
		balancePeriod.BalanceValidationR1 = status
	case measures.R2:
		balancePeriod.BalanceValidationR2 = status
	case measures.R3:
		balancePeriod.BalanceValidationR3 = status
	case measures.R4:
		balancePeriod.BalanceValidationR4 = status
	}
}

func (balancePeriod *BillingBalancePeriod) setPeriodMagnitude(magnitude measures.Magnitude, value float64) {
	switch magnitude {
	case measures.AI:
		balancePeriod.AI = value
	case measures.AE:
		balancePeriod.AE = value
	case measures.R1:
		balancePeriod.R1 = value
	case measures.R2:
		balancePeriod.R2 = value
	case measures.R3:
		balancePeriod.R3 = value
	case measures.R4:
		balancePeriod.R4 = value
	}
}

func (balancePeriod *BillingBalancePeriod) sumPeriodMagnitude(magnitude measures.Magnitude, value float64) {
	switch magnitude {
	case measures.AI:
		balancePeriod.AI += value
	case measures.AE:
		balancePeriod.AE += value
	case measures.R1:
		balancePeriod.R1 += value
	case measures.R2:
		balancePeriod.R2 += value
	case measures.R3:
		balancePeriod.R3 += value
	case measures.R4:
		balancePeriod.R4 += value
	}
}

func (balancePeriod *BillingBalancePeriod) getBalanceGeneralType(magnitude measures.Magnitude) GeneralEstimateMethod {
	switch magnitude {
	case measures.AI:
		return balancePeriod.BalanceGeneralTypeAI
	case measures.AE:
		return balancePeriod.BalanceGeneralTypeAE
	case measures.R1:
		return balancePeriod.BalanceGeneralTypeR1
	case measures.R2:
		return balancePeriod.BalanceGeneralTypeR2
	case measures.R3:
		return balancePeriod.BalanceGeneralTypeR3
	case measures.R4:
		return balancePeriod.BalanceGeneralTypeR4
	default:
		return ""
	}
}

func (balancePeriod *BillingBalancePeriod) setBalanceGeneralType(magnitude measures.Magnitude, generalType GeneralEstimateMethod) {
	switch magnitude {
	case measures.AI:
		balancePeriod.BalanceGeneralTypeAI = generalType
	case measures.AE:
		balancePeriod.BalanceGeneralTypeAE = generalType
	case measures.R1:
		balancePeriod.BalanceGeneralTypeR1 = generalType
	case measures.R2:
		balancePeriod.BalanceGeneralTypeR2 = generalType
	case measures.R3:
		balancePeriod.BalanceGeneralTypeR3 = generalType
	case measures.R4:
		balancePeriod.BalanceGeneralTypeR4 = generalType
	}
}

func (balancePeriod BillingBalancePeriod) Values() *measures.Values {
	return &measures.Values{
		AI: balancePeriod.AI,
		AE: balancePeriod.AE,
		R1: balancePeriod.R1,
		R2: balancePeriod.R2,
		R3: balancePeriod.R3,
		R4: balancePeriod.R4,
	}
}
