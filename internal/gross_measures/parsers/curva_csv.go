package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CurvaCsv struct {
	Parser
}

func NewCurvaCsv(parser Parser) CurvaCsv {
	return CurvaCsv{
		parser,
	}
}

func (c CurvaCsv) decimalToBinary(s ...string) (string, error) {

	sumDecimal := int64(0)

	for _, v := range s {
		intBin, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return "", err
		}
		sumDecimal += intBin
	}

	stringBin := strconv.FormatInt(sumDecimal, 2)
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

func (c CurvaCsv) formatRow(row []string) (gross_measures.MeasureCurveWrite, error) {

	if len(row) != 16 {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	endSeason := row[3]
	endDate, err := c.FormatDateBySeason("02/01/06 15:04", row[2], rune(endSeason[0]))

	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	ai, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	AE, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	R1, err := strconv.ParseFloat(row[8], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	R2, err := strconv.ParseFloat(row[10], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	R3, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	R4, err := strconv.ParseFloat(row[14], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	qualifier, err := c.decimalToBinary(row[5], row[7], row[9], row[11], row[13], row[15])
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	return gross_measures.MeasureCurveWrite{
		EndDate:           endDate,
		ReadingDate:       time.Time{},
		Type:              measures.Incremental,
		ReadingType:       measures.Curve,
		CurveType:         measures.HourlyMeasureCurveReadingType,
		MeterSerialNumber: row[0],
		ConcentratorID:    "",
		File:              "",
		DistributorID:     "",
		Origin:            measures.File,
		Qualifier:         qualifier,
		AI:                ai * 1000,
		AE:                AE * 1000,
		R1:                R1 * 1000,
		R2:                R2 * 1000,
		R3:                R3 * 1000,
		R4:                R4 * 1000,
	}, nil
}

func (c CurvaCsv) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCurveWrite, error) {
	storage, err := c.storage(ctx)
	if err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}

	defer storage.Close()

	reader, err := storage.NewReader(ctx, dto.FilePath)

	if err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	isFirstLine := true
	measureList := make([]gross_measures.MeasureCurveWrite, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []gross_measures.MeasureCurveWrite{}, err
		}
		if isFirstLine {
			isFirstLine = false
			continue
		}

		distributorCdos := c.GetDistributorFromFilePath(dto.FilePath)

		measure, err := c.formatRow(row)

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

		measureList = append(measureList, measure)
	}

	return measureList, nil
}

func (c CurvaCsv) getMeasureDateFromFilePath(filePath string) (time.Time, error) {
	filename := c.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 4 {
		return time.Time{}, errors.New("invalid date")
	}
	return c.FormatDate("0601020304", filenamePaths[2])
}

func (c CurvaCsv) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

	filename := c.GetFileNameFromFilePath(dto.FilePath)

	if filename == "" {
		return false
	}

	ext := filepath.Ext(filename)

	if ext != ".csv" {
		return false
	}

	filenamePaths := strings.Split(filename, "_")

	if len(filenamePaths) != 4 {
		return false
	}

	return strings.Contains(filename, "CCCH")
}
