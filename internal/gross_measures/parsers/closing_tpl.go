package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ClosingTpl struct {
	Parser
	ClosingFileName string
}

func NewClosingTpl(parser Parser) ClosingTpl {
	return ClosingTpl{
		Parser: parser,
	}
}

type ClosingIndexStruct struct {
	Contract    int
	Period      int
	StartDate   int
	InitStation int
	EndDate     int
	EndStation  int
	AAbs        int
	AInc        int
	S1          int
	RiAbs       int
	RiInc       int
	S2          int
	RcAbs       int
	RcInc       int
	S3          int
	ExcA        int
	S4          int
	MaxA        int
	S5          int
	MaxADate    int
}

var closingConfig = ClosingIndexStruct{
	Contract:    0,
	Period:      1,
	StartDate:   2,
	InitStation: 3,
	EndDate:     4,
	EndStation:  5,
	AAbs:        6,
	AInc:        7,
	S1:          8,
	RiAbs:       9,
	RiInc:       10,
	S2:          11,
	RcAbs:       12,
	RcInc:       13,
	S3:          14,
	ExcA:        15,
	S4:          16,
	MaxA:        17,
	S5:          18,
	MaxADate:    19,
}

type ValueContractConfig struct {
	ContractIncremental   string
	IndexValueAbsolute    int
	IndexValueIncremental int
	OnlyMatchContract     bool
}

type ValueConfig struct {
	AI ValueContractConfig
	AE ValueContractConfig
	R1 ValueContractConfig
	R2 ValueContractConfig
	R3 ValueContractConfig
	R4 ValueContractConfig
	Mx ValueContractConfig
	Fx ValueContractConfig
	E  ValueContractConfig
}

var valueConfig = ValueConfig{
	AI: ValueContractConfig{
		ContractIncremental:   "1",
		IndexValueAbsolute:    closingConfig.AAbs,
		IndexValueIncremental: closingConfig.AInc,
		OnlyMatchContract:     false,
	},

	AE: ValueContractConfig{
		ContractIncremental:   "3",
		IndexValueAbsolute:    closingConfig.AAbs,
		IndexValueIncremental: closingConfig.AInc,
		OnlyMatchContract:     false,
	},

	R1: ValueContractConfig{
		ContractIncremental:   "1",
		IndexValueAbsolute:    closingConfig.RiAbs,
		IndexValueIncremental: closingConfig.RiInc,
		OnlyMatchContract:     false,
	},

	R2: ValueContractConfig{
		ContractIncremental:   "1",
		IndexValueAbsolute:    closingConfig.RcAbs,
		IndexValueIncremental: closingConfig.RcInc,
		OnlyMatchContract:     false,
	},

	R3: ValueContractConfig{
		ContractIncremental:   "3",
		IndexValueAbsolute:    closingConfig.RiAbs,
		IndexValueIncremental: closingConfig.RiInc,
		OnlyMatchContract:     false,
	},

	R4: ValueContractConfig{
		ContractIncremental:   "3",
		IndexValueAbsolute:    closingConfig.RcAbs,
		IndexValueIncremental: closingConfig.RcInc,
		OnlyMatchContract:     false,
	},

	Mx: ValueContractConfig{
		ContractIncremental:   "1",
		IndexValueAbsolute:    closingConfig.MaxA,
		IndexValueIncremental: closingConfig.MaxA,
		OnlyMatchContract:     true,
	},

	Fx: ValueContractConfig{
		ContractIncremental:   "1",
		IndexValueAbsolute:    closingConfig.MaxADate,
		IndexValueIncremental: closingConfig.MaxADate,
		OnlyMatchContract:     true,
	},

	E: ValueContractConfig{
		ContractIncremental:   "1",
		IndexValueAbsolute:    closingConfig.ExcA,
		IndexValueIncremental: closingConfig.ExcA,
		OnlyMatchContract:     true,
	},
}

func (r ClosingTpl) getValueIndex(config ValueContractConfig, contract string, typeValue measures.Type, row []string) (value string) {
	if config.OnlyMatchContract == true {
		if config.ContractIncremental == contract {
			return row[config.IndexValueIncremental]
		}
		return ""
	}

	if typeValue == measures.Absolute {
		return row[config.IndexValueAbsolute]
	} else {
		if config.ContractIncremental == contract {
			return row[config.IndexValueIncremental]
		}
		return ""
	}
}

func (c ClosingTpl) formatRow(row []string, typeRow measures.Type) (gross_measures.MeasureCloseWrite, error) {

	if len(row) != 20 {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure data")
	}

	startSeason := row[closingConfig.EndStation]
	startDate, err := c.formatDate("02/01/06 03:04", row[closingConfig.EndDate], rune(startSeason[0]))
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure start date")
	}

	endSeason := row[closingConfig.InitStation]
	endDate, err := c.formatDate("02/01/06 03:04", row[closingConfig.StartDate], rune(endSeason[0]))
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure end date")
	}

	contract := row[closingConfig.Contract]
	period := c.FormatPeriod(row[closingConfig.Period])

	measureClosePeriod := gross_measures.MeasureClosePeriod{
		Period: period,
	}

	rowValueAI := c.getValueIndex(valueConfig.AI, contract, typeRow, row)
	if rowValueAI != "" {
		AI, err := strconv.ParseFloat(rowValueAI, 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure AI")
		}
		measureClosePeriod.AI = AI
	}

	rowValueAE := c.getValueIndex(valueConfig.AE, contract, typeRow, row)
	if rowValueAE != "" {
		AE, err := strconv.ParseFloat(rowValueAE, 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure AE")
		}
		measureClosePeriod.AE = AE
	}

	rowValueR1 := c.getValueIndex(valueConfig.R1, contract, typeRow, row)
	if rowValueR1 != "" {
		R1, err := strconv.ParseFloat(rowValueR1, 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R1")
		}
		measureClosePeriod.R1 = R1
	}

	rowValueR2 := c.getValueIndex(valueConfig.R2, contract, typeRow, row)
	if rowValueR2 != "" {
		R2, err := strconv.ParseFloat(rowValueR2, 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R2")
		}
		measureClosePeriod.R2 = R2
	}

	rowValueR3 := c.getValueIndex(valueConfig.R3, contract, typeRow, row)
	if rowValueR3 != "" {
		R3, err := strconv.ParseFloat(rowValueR3, 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R3")
		}
		measureClosePeriod.R3 = R3
	}

	rowValueR4 := c.getValueIndex(valueConfig.R4, contract, typeRow, row)
	if rowValueR4 != "" {
		R4, err := strconv.ParseFloat(rowValueR4, 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R4")
		}
		measureClosePeriod.R4 = R4
	}

	if contract == "1" {
		FX, err := time.Parse("02/01/06 03:04", c.getValueIndex(valueConfig.Fx, contract, typeRow, row))
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure FX")
		}

		MX, err := strconv.ParseFloat(c.getValueIndex(valueConfig.Mx, contract, typeRow, row), 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure Mx")
		}

		E, err := strconv.ParseFloat(c.getValueIndex(valueConfig.E, contract, typeRow, row), 64)
		if err != nil {
			return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure E")
		}

		measureClosePeriod.FX = FX
		measureClosePeriod.MX = MX
		measureClosePeriod.E = E
	}

	newMeasure := gross_measures.MeasureCloseWrite{
		StartDate:         startDate,
		EndDate:           endDate,
		ReadingDate:       time.Time{},
		Type:              typeRow,
		ReadingType:       measures.BillingClosure,
		Contract:          contract,
		Periods:           []gross_measures.MeasureClosePeriod{measureClosePeriod},
		MeterSerialNumber: "",
		Origin:            measures.TPL,
	}

	return newMeasure, nil
}

func (c ClosingTpl) formatDate(layout, date string, season rune) (time.Time, error) {
	return c.FormatDateBySeason(layout, date, season)
}

func (c ClosingTpl) getMeasureDateFromFilePath(filePath string) (time.Time, error) {
	filename := c.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return time.Time{}, errors.New("invalid date")
	}
	m1 := regexp.MustCompile(`\D+`)
	date := m1.ReplaceAllString(filenamePaths[4], "")
	if len(date) != 8 {
		return time.Time{}, errors.New("invalid date")
	}

	return c.FormatDate("20060102", date)
}

func (p ClosingTpl) getFileNameFromFilePath(filePath string) string {
	paths := strings.Split(filePath, "/")
	if len(paths) != 5 {
		return ""
	}
	return paths[4]
}

func (p ClosingTpl) getEquipmentId(filePath string) string {
	filename := p.getFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return ""
	}

	return filenamePaths[1]
}

func (c ClosingTpl) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error) {
	storage, err := c.storage(ctx)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	defer storage.Close()

	reader, err := storage.NewReader(ctx, dto.FilePath)

	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '|'

	measuresParsed := make(map[string]gross_measures.MeasureCloseWrite, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []gross_measures.MeasureCloseWrite{}, err
		}
		contract := row[closingConfig.Contract]

		if contract != "1" && contract != "3" {
			continue
		}
		distributorCdos := c.GetDistributorFromFilePath(dto.FilePath)

		typesMeasures := []measures.Type{measures.Absolute, measures.Incremental}

		for _, typeM := range typesMeasures {
			measure, err := c.formatRow(row, typeM)

			if err != nil {
				continue
			}

			MeasureDate, err := c.getMeasureDateFromFilePath(dto.FilePath)

			if err != nil {
				continue
			}

			measure.File = dto.FilePath
			measure.DistributorCDOS = distributorCdos
			measure.ReadingDate = MeasureDate.UTC()
			measure.MeterSerialNumber = c.getEquipmentId(dto.FilePath)

			key := fmt.Sprintf("%s%v%v%s", measure.Type, measure.EndDate, measure.StartDate, measure.Contract)

			savedMeasure, ok := measuresParsed[key]

			if ok {
				measure.Periods = append(savedMeasure.Periods, measure.Periods...)
			}
			measuresParsed[key] = measure
		}
	}

	return utils.MapToSlice(measuresParsed), nil
}

func (c ClosingTpl) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

	filename := c.GetFileNameFromFilePath(dto.FilePath)

	if filename == "" {
		return false
	}

	ext := filepath.Ext(filename)

	if ext != ".tpl" {
		return false
	}

	filenamePaths := strings.Split(filename, "_")

	if len(filenamePaths) != 5 {
		return false
	}

	return strings.Contains(filename, "iivm")
}
