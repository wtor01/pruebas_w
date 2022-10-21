package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

type CloseCsv struct {
	Parser
}

type rowPositions struct {
	ai  int
	AE  int
	R1  int
	Ri1 int
	R2  int
	R3  int
	R4  int
	Rc4 int
	M1  int
	F1  int
	E1  int
}

var period1 = rowPositions{
	ai:  3,
	AE:  63,
	R1:  39,
	Ri1: 27,
	R2:  45,
	R3:  51,
	R4:  57,
	Rc4: 33,
	M1:  9,
	F1:  10,
	E1:  21,
}
var period2 = rowPositions{
	ai:  4,
	AE:  64,
	R1:  40,
	Ri1: 28,
	R2:  46,
	R3:  52,
	R4:  58,
	Rc4: 34,
	M1:  11,
	F1:  12,
	E1:  22,
}
var period3 = rowPositions{
	ai:  5,
	AE:  65,
	R1:  41,
	Ri1: 29,
	R2:  47,
	R3:  53,
	R4:  59,
	Rc4: 35,
	M1:  13,
	F1:  14,
	E1:  23,
}
var period4 = rowPositions{
	ai:  6,
	AE:  66,
	R1:  42,
	Ri1: 30,
	R2:  48,
	R3:  54,
	R4:  60,
	Rc4: 36,
	M1:  15,
	F1:  16,
	E1:  24,
}
var period5 = rowPositions{
	ai:  7,
	AE:  67,
	R1:  43,
	Ri1: 31,
	R2:  49,
	R3:  55,
	R4:  61,
	Rc4: 37,
	M1:  17,
	F1:  18,
	E1:  25,
}
var period6 = rowPositions{
	ai:  8,
	AE:  68,
	R1:  44,
	Ri1: 32,
	R2:  50,
	R3:  56,
	R4:  62,
	Rc4: 38,
	M1:  19,
	F1:  20,
	E1:  26,
}

func NewCloseCsv(parser Parser) CloseCsv {
	return CloseCsv{
		parser,
	}
}

func (c CloseCsv) getRowsPositions(index int) rowPositions {
	switch index {
	case 0:
		return period1
	case 1:
		return period2
	case 2:
		return period3
	case 3:
		return period4
	case 4:
		return period5
	case 5:
		return period6
	}

	return rowPositions{}
}

func (c CloseCsv) parseRowFloat(row []string, index int) float64 {
	value, err := strconv.ParseFloat(row[index], 64)
	if err != nil {
		return 0
	}
	value = value * 1000

	return value
}

func (c CloseCsv) formatRow(row []string, filePath string) (gross_measures.MeasureCloseWrite, error) {
	distributorCdos := c.GetDistributorFromFilePath(filePath)

	if len(row) != 69 {
		return gross_measures.MeasureCloseWrite{}, gross_measures.ErrInvalidMeasureData
	}

	startDate, err := c.FormatDate("02/01/2006", row[1])
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, gross_measures.ErrInvalidMeasureDate
	}
	measureDate, err := c.FormatDate("02/01/2006", row[2])
	if err != nil {
		return gross_measures.MeasureCloseWrite{}, gross_measures.ErrInvalidMeasureDate
	}

	periods := make([]gross_measures.MeasureClosePeriod, 0, 6)

	for period := 0; period < (gross_measures.MeasureClose{}).NumberPeriods(); period++ {
		positions := c.getRowsPositions(period)
		if positions.ai == 0 {
			continue
		}
		ai := c.parseRowFloat(row, positions.ai)
		AE := c.parseRowFloat(row, positions.AE)
		R1 := c.parseRowFloat(row, positions.R1)
		if R1 == 0 {
			R1 = c.parseRowFloat(row, positions.Ri1)
		}
		R2 := c.parseRowFloat(row, positions.R2)
		R3 := c.parseRowFloat(row, positions.R3)
		R4 := c.parseRowFloat(row, positions.R4)
		if R4 == 0 {
			R4 = c.parseRowFloat(row, positions.Rc4)
		}
		M1 := c.parseRowFloat(row, positions.M1)
		F1, _ := c.FormatDate("02/01/2006 15:04", row[positions.F1])

		E1 := c.parseRowFloat(row, positions.E1)
		P := c.FormatPeriod(fmt.Sprintf("%d", period+1))
		periods = append(periods, gross_measures.MeasureClosePeriod{
			Period: P,
			AI:     ai,
			AE:     AE,
			R1:     R1,
			R2:     R2,
			R3:     R3,
			R4:     R4,
			MX:     M1,
			FX:     F1.UTC(),
			E:      E1,
		})
	}

	measure := gross_measures.MeasureCloseWrite{
		StartDate:         startDate.UTC(),
		EndDate:           measureDate.UTC(),
		ReadingDate:       measureDate.UTC(),
		Type:              measures.Absolute,
		ReadingType:       measures.BillingClosure,
		MeterSerialNumber: row[0],
		ConcentratorID:    "",
		File:              filePath,
		DistributorID:     "",
		DistributorCDOS:   distributorCdos,
		Periods:           periods,
		Origin:            measures.STM,
	}

	return measure, nil
}

func (c CloseCsv) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error) {
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
	measureList := make([]gross_measures.MeasureCloseWrite, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []gross_measures.MeasureCloseWrite{}, err
		}

		if isFirstLine {
			isFirstLine = false
			continue
		}

		mm, _ := c.formatRow(row, dto.FilePath)

		measureList = append(measureList, mm)

	}

	return measureList, nil
}

func (c CloseCsv) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

	if dto.FilePath == "" {
		return false
	}

	filename := c.GetFileNameFromFilePath(dto.FilePath)

	ext := filepath.Ext(filename)

	if ext != ".csv" {
		return false
	}

	filenamePaths := strings.Split(filename, "_")

	if len(filenamePaths) == 6 || len(filenamePaths) == 5 {
		fmt.Print(filenamePaths)
		return strings.Contains(filename, "CIERRE")
	}
	return false

}
