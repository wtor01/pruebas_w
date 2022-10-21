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

type CurveTpl struct {
	Parser
}

func NewCurveTpl(parser Parser) CurveTpl {
	return CurveTpl{
		parser,
	}
}

func (c CurveTpl) decimalToBinary(s ...string) (string, error) {

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

func (c CurveTpl) formatDate(layout, date string, season rune) (time.Time, error) {
	return c.FormatDateBySeason(layout, date, season)
}

func (c CurveTpl) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {
	if !strings.Contains(dto.FilePath, "TPL") {
		return false
	}

	filename := c.GetFileNameFromFilePath(dto.FilePath)

	ext := filepath.Ext(filename)

	if ext != ".tpl" {
		return false
	}

	if !strings.Contains(filename, "iicc") {
		return false
	}

	return true
}

func (c CurveTpl) getEquipmentIDFromFilePath(filePath string) string {
	filename := c.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return ""
	}
	return filenamePaths[1]
}

func (c CurveTpl) getMeasureDateFromFilePath(filePath string) (time.Time, error) {
	filename := c.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return time.Time{}, errors.New("invalid date")
	}
	filenamePathsDate := strings.Split(filenamePaths[4], ".")
	return c.FormatDate("20060102", filenamePathsDate[0])
}

func (c CurveTpl) formatRow(row []string) (gross_measures.MeasureCurveWrite, error) {

	if len(row) != 16 {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure data")
	}

	typ := measures.Absolute

	if row[1] == "1" {
		typ = measures.Incremental
	}

	endSeason := row[3]
	endDate, err := c.formatDate("02/01/06 15:04", row[2], rune(endSeason[0]))

	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure end_date")
	}

	AI, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure AI")
	}

	AE, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure AE")
	}

	R1, err := strconv.ParseFloat(row[8], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure R1")
	}

	R2, err := strconv.ParseFloat(row[10], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure R2")
	}

	R3, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure R3")
	}

	R4, err := strconv.ParseFloat(row[14], 64)
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure R4")
	}

	qualifier, err := c.decimalToBinary(row[5], row[7], row[9], row[11], row[13], row[15])
	if err != nil {
		return gross_measures.MeasureCurveWrite{}, errors.New("invalid measure qualifier")
	}

	return gross_measures.MeasureCurveWrite{
		EndDate:           endDate,
		ReadingDate:       time.Time{},
		Type:              typ,
		ReadingType:       measures.Curve,
		CurveType:         measures.HourlyMeasureCurveReadingType,
		MeterSerialNumber: "",
		File:              "",
		DistributorID:     "",
		Origin:            measures.TPL,
		Qualifier:         qualifier,
		AI:                AI * 1000,
		AE:                AE * 1000,
		R1:                R1 * 1000,
		R2:                R2 * 1000,
		R3:                R3 * 1000,
		R4:                R4 * 1000,
	}, nil
}

func (c CurveTpl) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCurveWrite, error) {
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
	csvReader.Comma = '|'
	measureList := make([]gross_measures.MeasureCurveWrite, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []gross_measures.MeasureCurveWrite{}, err
		}

		distributorCdos := c.GetDistributorFromFilePath(dto.FilePath)

		equipmentID := c.getEquipmentIDFromFilePath(dto.FilePath)

		if equipmentID == "" {
			continue
		}

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
		measure.MeterSerialNumber = equipmentID
		measure.ReadingDate = MeasureDate.UTC()

		measureList = append(measureList, measure)
	}

	return measureList, nil
}
