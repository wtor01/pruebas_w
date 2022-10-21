package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ReadingsTpl struct {
	Parser
	HeaderFileName   string
	ReadingsFileName string
}

func NewReadingsTpl(parser Parser) ReadingsTpl {
	return ReadingsTpl{
		Parser:           parser,
		HeaderFileName:   "CABL.TPL",
		ReadingsFileName: "LEC.TPL",
	}
}

func (r ReadingsTpl) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {
	if !strings.Contains(dto.FilePath, "TPL") {
		return false
	}

	filename := r.GetFileNameFromFilePath(dto.FilePath)

	ext := filepath.Ext(filename)

	if ext != ".zip" {
		return false
	}

	if strings.Contains(filename, "iivm") || strings.Contains(filename, "iicc") {
		return false
	}

	return true
}

func (r ReadingsTpl) generateMeasuresFromHeaderFile(csvReader *csv.Reader, filePath string) (map[string]gross_measures.MeasureCloseWrite, error) {
	measureList := make(map[string]gross_measures.MeasureCloseWrite)

	rowsIdsToDelete := make([]string, 0)
	distributorCdos := r.GetDistributorFromFilePath(filePath)
	endDate, err := r.GetMeasureDateFromFilePath(filePath, 0)
	if err != nil {
		return measureList, err
	}
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return measureList, err
		}
		if len(row) != 4 {
			continue
		}
		measureDate, err := r.FormatDate("20060102", row[2])
		if err != nil {
			continue
		}

		measureDate = measureDate.UTC().Add(time.Hour * time.Duration(24-measureDate.UTC().Hour()))
		rowsIdsToDelete = append(rowsIdsToDelete, row[0])

		measureList[row[0]] = gross_measures.MeasureCloseWrite{
			EndDate:           endDate.UTC(),
			ReadingDate:       measureDate,
			Type:              measures.Absolute,
			ReadingType:       measures.DailyClosure,
			MeterSerialNumber: row[1],
			ConcentratorID:    "",
			File:              filePath,
			DistributorCDOS:   distributorCdos,
			Origin:            measures.TPL,
			Qualifier:         "",
			Periods:           make([]gross_measures.MeasureClosePeriod, 0),
		}
	}

	return measureList, nil
}

func (r ReadingsTpl) setMeasuresValues(csvReader *csv.Reader, measuresMap map[string]gross_measures.MeasureCloseWrite) error {
	measurePeriod := make(map[string]gross_measures.MeasureClosePeriod, 0)
	separator := "|||"
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(row) != 4 {
			continue
		}

		rowID := row[0]
		period := row[2]
		periodId := fmt.Sprintf("%s%s%s", rowID, separator, period)
		magnitude := row[1]
		value := row[3]

		p, _ := measurePeriod[periodId]

		p.Period = measures.PeriodKey(period)

		err = r.setRowValuesToMeasure(&p, magnitude, value)
		measurePeriod[periodId] = p
	}

	for id, period := range measurePeriod {
		ids := strings.Split(id, separator)
		rowId := ids[0]
		m, ok := measuresMap[rowId]
		if !ok {
			continue
		}

		m.Periods = append(m.Periods, period)
		measuresMap[rowId] = m
	}

	return nil
}

func (r ReadingsTpl) getFileReaders(ctx context.Context, dto gross_measures.HandleFileDTO) (headerReader *csv.Reader, bodyReader *csv.Reader, err error) {
	storage, err := r.storage(ctx)
	if err != nil {
		return headerReader, bodyReader, err
	}

	defer storage.Close()

	reader, err := storage.NewReader(ctx, dto.FilePath)
	tmpFileZip := r.GetZipTmpFilePath(dto.FilePath)
	tmpFile := r.GetTmpFilePath(dto.FilePath)

	f, err := os.Create(tmpFileZip)

	if err != nil {
		return headerReader, bodyReader, err
	}

	defer func(f *os.File) {
		f.Close()
		err = os.Remove(tmpFileZip)
	}(f)

	_, err = io.Copy(f, reader)
	if err != nil {
		return headerReader, bodyReader, err
	}
	err = r.Unzip(tmpFileZip, tmpFile)

	if err != nil {
		return headerReader, bodyReader, err
	}

	defer func() {
		_ = os.RemoveAll(tmpFile)
	}()

	headerContent, err := os.ReadFile(filepath.Join(tmpFile, r.HeaderFileName))
	if err != nil {
		return headerReader, bodyReader, err
	}

	headerReader = csv.NewReader(bytes.NewReader(headerContent))
	headerReader.Comma = '|'

	readingsContent, err := os.ReadFile(filepath.Join(tmpFile, r.ReadingsFileName))
	if err != nil {
		return headerReader, bodyReader, err
	}
	bodyReader = csv.NewReader(bytes.NewReader(readingsContent))
	bodyReader.Comma = '|'

	return headerReader, bodyReader, nil
}

func (r ReadingsTpl) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error) {

	headerReader, bodyReader, err := r.getFileReaders(ctx, dto)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	measureList, err := r.generateMeasuresFromHeaderFile(headerReader, dto.FilePath)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	err = r.setMeasuresValues(bodyReader, measureList)

	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	mm := make([]gross_measures.MeasureCloseWrite, 0, len(measureList))

	for _, m := range measureList {
		if len(m.Periods) == 0 {
			continue
		}
		mm = append(mm, m)
	}

	return mm, nil
}

func (r ReadingsTpl) setRowValuesToMeasure(measure *gross_measures.MeasureClosePeriod, magnitude, value string) error {
	switch magnitude {
	case "AbsA+":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.AI = num * 1000
			break
		}
	case "AbsA-":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.AE = num * 1000
			break
		}
	case "AbsRi+":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.R1 = num * 1000
			break
		}
	case "AbsRc+":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.R2 = num * 1000
			break
		}
	case "AbsRi-":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.R3 = num * 1000
			break
		}
	case "AbsRc-":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.R4 = num * 1000
			break
		}
	case "MaxA":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.MX = num * 1000
			break
		}
	case "ExcA":
		{
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			measure.E = num * 1000
			break
		}
	default:
		break
	}
	return nil
}
