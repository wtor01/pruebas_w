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
	"strconv"
	"strings"
)

type ResumenCsv struct {
	Parser
}

func NewResumenCsv(parser Parser) ResumenCsv {
	return ResumenCsv{
		parser,
	}
}

type PeriodIndexConfig struct {
	AI int
	AE int
	R1 int
	R2 int
	R3 int
	R4 int
}

var periodConfig0 = PeriodIndexConfig{
	AI: 5,
	AE: 12,
	R1: 19,
	R2: 26,
	R3: 33,
	R4: 40,
}

var periodConfig1 = PeriodIndexConfig{
	AI: 6,
	AE: 13,
	R1: 20,
	R2: 27,
	R3: 34,
	R4: 41,
}
var periodConfig2 = PeriodIndexConfig{
	AI: 7,
	AE: 14,
	R1: 21,
	R2: 28,
	R3: 35,
	R4: 42,
}
var periodConfig3 = PeriodIndexConfig{
	AI: 8,
	AE: 15,
	R1: 22,
	R2: 29,
	R3: 36,
	R4: 43,
}
var periodConfig4 = PeriodIndexConfig{
	AI: 9,
	AE: 16,
	R1: 23,
	R2: 30,
	R3: 37,
	R4: 44,
}
var periodConfig5 = PeriodIndexConfig{
	AI: 10,
	AE: 17,
	R1: 24,
	R2: 31,
	R3: 38,
	R4: 45,
}
var periodConfig6 = PeriodIndexConfig{
	AI: 11,
	AE: 18,
	R1: 25,
	R2: 32,
	R3: 39,
	R4: 46,
}

var periodsConfig = map[measures.PeriodKey]PeriodIndexConfig{
	measures.P0: periodConfig0,
	measures.P1: periodConfig1,
	measures.P2: periodConfig2,
	measures.P3: periodConfig3,
	measures.P4: periodConfig4,
	measures.P5: periodConfig5,
	measures.P6: periodConfig6,
}

func (c ResumenCsv) cleanValue(value string) string {
	return strings.ReplaceAll(string(value), "\u0000", "")
}

func (c ResumenCsv) formatRow(row []string, period measures.PeriodKey, periodConfig PeriodIndexConfig) (gross_measures.MeasureCloseWrite, error) {

	if len(row) != 48 {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure data")
	}

	measureDate, err := c.FormatDate("02/01/2006 03:04:05", c.cleanValue(row[3]))

	endSeason := row[2]
	endDate, err := c.FormatDateBySeason("02/01/2006 3:04:05", c.cleanValue(row[1]), rune(endSeason[0]))

	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure date")
	}

	contract := c.cleanValue(row[4])

	AI, err := strconv.ParseFloat(c.cleanValue(row[periodConfig.AI]), 64)
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure AI")
	}

	AE, err := strconv.ParseFloat(c.cleanValue(row[periodConfig.AE]), 64)
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure AE")
	}

	R1, err := strconv.ParseFloat(c.cleanValue(row[periodConfig.R1]), 64)
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R1")
	}

	R2, err := strconv.ParseFloat(c.cleanValue(row[periodConfig.R2]), 64)
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R2")
	}

	R3, err := strconv.ParseFloat(c.cleanValue(row[periodConfig.R3]), 64)
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R3")
	}

	R4, err := strconv.ParseFloat(c.cleanValue(row[periodConfig.R4]), 64)
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, errors.New("invalid measure R4")
	}

	newMeasure := gross_measures.MeasureCloseWrite{
		EndDate:           endDate,
		ReadingDate:       measureDate.UTC(),
		Type:              measures.Absolute,
		ReadingType:       measures.DailyClosure,
		Contract:          contract,
		MeterSerialNumber: c.cleanValue(row[0]),
		ConcentratorID:    "",
		File:              "",
		Origin:            measures.File,
		Periods: []gross_measures.MeasureClosePeriod{
			{
				Period: period,
				AI:     AI * 1000,
				AE:     AE * 1000,
				R1:     R1 * 1000,
				R2:     R2 * 1000,
				R3:     R3 * 1000,
				R4:     R4 * 1000,
			},
		},
	}
	return newMeasure, nil
}

func (c ResumenCsv) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error) {
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
	csvReader.Comma = ';'
	isFirstLine := true
	measuresParsed := make(map[string]gross_measures.MeasureCloseWrite, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			continue
		}

		if isFirstLine {
			isFirstLine = false
			continue
		}
		for periodNumber, configPeriod := range periodsConfig {

			distributorCdos := c.GetDistributorFromFilePath(dto.FilePath)

			measure, err := c.formatRow(row, periodNumber, configPeriod)

			if err != nil {
				continue
			}

			measure.File = dto.FilePath
			measure.DistributorCDOS = distributorCdos

			key := fmt.Sprintf("%v%v%s%s", measure.EndDate, measure.ReadingDate, measure.Contract, measure.MeterSerialNumber)

			savedMeasure, ok := measuresParsed[key]

			if ok {
				measure.Periods = append(savedMeasure.Periods, measure.Periods...)
			}
			measuresParsed[key] = measure
		}
	}
	return utils.MapToSlice(measuresParsed), nil
}

func (c ResumenCsv) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

	filename := c.GetFileNameFromFilePath(dto.FilePath)

	if filename == "" {
		return false
	}

	ext := filepath.Ext(filename)

	if ext != ".csv" {
		return false
	}

	filenamePaths := strings.Split(filename, "_")

	if len(filenamePaths) != 5 {
		return false
	}

	return strings.Contains(filename, "S05")
}
