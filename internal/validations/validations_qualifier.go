package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"strconv"
	"strings"
)

/*
solo se aplican a curva

*/

var binariesInvalidQualifier = [8]int{
	1,
	2,
	3,
	5,
	6,
	7,
}

type ValidatorQualifier struct {
	ValidatorBase
}

var validDateQualifierKeys = map[string]struct{}{
	Qualifier: {},
}

func NewValidatorQualifier(v ValidationData, status measures.Status) (ValidatorQualifier, error) {
	if len(v.Config) != 0 {
		return ValidatorQualifier{}, newErrorConfigProperties()
	}

	validator := ValidatorQualifier{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validDateQualifierKeys); err != nil {
		return ValidatorQualifier{}, err
	}

	return validator, nil
}

func (v ValidatorQualifier) decimalToBinary(s string) (string, error) {
	if len(s) == 8 {
		return s, nil
	}
	intBin, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", err
	}
	stringBin := strconv.FormatInt(intBin, 2)
	if len(stringBin) != 8 && len(stringBin) < 8 {
		toAdd := make([]int, 8-len(stringBin))
		var sb strings.Builder
		for range toAdd {
			sb.WriteString("0")
		}
		sb.WriteString(stringBin)
		stringBin = sb.String()
	}
	return stringBin, nil
}

func (v ValidatorQualifier) reverseBinary(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func (v ValidatorQualifier) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType || !m.isInWhiteList(Qualifier) {
		return nil
	}
	for _, k := range v.Keys {
		if k != Qualifier {
			continue
		}

		qualifierBinary, err := v.decimalToBinary(m.Qualifier)
		qualifierReversed := v.reverseBinary(qualifierBinary)

		if err != nil {
			return &v.ValidatorBase
		}
		for _, binInv := range binariesInvalidQualifier {
			if string(qualifierReversed[binInv]) == "1" {
				return &v.ValidatorBase
			}
		}
	}
	return nil
}
