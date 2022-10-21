package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type PrimeS05 struct {
	Parser
}

type ReportS05 struct {
	XMLName xml.Name `xml:"Report"`
	IdRpt   string   `xml:"IdRpt,attr"`
	CncS05  CncS05   `xml:"Cnc"`
}

type CncS05 struct {
	XMLName xml.Name `xml:"Cnc"`
	Id      string   `xml:"Id,attr"`
	CntS05  []CntS05 `xml:"Cnt"`
}

type CntS05 struct {
	XMLName xml.Name `xml:"Cnt"`
	Id      string   `xml:"Id,attr"`
	S05     []S05    `xml:"S05"`
	ErrCat  string   `xml:"ErrCat,attr"`
	ErrCode string   `xml:"ErrCode,attr"`
}

type S05 struct {
	XMLName xml.Name `xml:"S05"`
	Fh      string   `xml:"Fh,attr"`
	Ctr     string   `xml:"Ctr,attr"`
	Pt      string   `xml:"Pt,attr"`
	Value   Value    `xml:"Value"`
}

type Value struct {
	XMLName xml.Name `xml:"Value"`
	AI      float64  `xml:"AIa,attr"`
	AE      float64  `xml:"AEa,attr"`
	R1      float64  `xml:"R1a,attr"`
	R2      float64  `xml:"R2a,attr"`
	R3      float64  `xml:"R3a,attr"`
	R4      float64  `xml:"R4a,attr"`
}

func NewPrimeS05(parser Parser) PrimeS05 {
	return PrimeS05{
		parser,
	}
}

func (p PrimeS05) formatDate(date string) (time.Time, error) {
	m1 := regexp.MustCompile(`\D+`)
	f := m1.ReplaceAllString(date, "")
	if len(f) < 14 {
		return time.Time{}, errors.New("invalid date")
	}
	milliseconds := f[14:]
	layout := "20060102150405" + strings.Repeat("0", utf8.RuneCountInString(milliseconds))

	season := date[len(date)-1]
	parsedDate, err := p.FormatDateBySeason(layout, f, rune(season))

	if err != nil {
		parsedDate, err = p.FormatDate(layout, f)
	}

	return parsedDate, err
}

func (p PrimeS05) getMeasureDateFromFilePath(filePath string) string {
	filename := p.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return ""
	}

	return filenamePaths[4]
}

func (p PrimeS05) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

	filename := p.GetFileNameFromFilePath(dto.FilePath)

	if filename == "" {
		return false
	}

	if !strings.Contains(dto.FilePath, "Prime") {
		return false
	}

	filenamePaths := strings.Split(filename, "_")

	if len(filenamePaths) != 5 {
		return false
	}

	return filenamePaths[2] == "S05"
}

func (p PrimeS05) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error) {
	storage, err := p.storage(ctx)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	defer storage.Close()

	content, err := storage.ReadAll(ctx, dto.FilePath)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}
	var report ReportS05
	err = xml.Unmarshal(content, &report)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	measuresParsed := make(map[string]gross_measures.MeasureCloseWrite, 0)

	for _, cnt := range report.CncS05.CntS05 {
		if cnt.ErrCode != "" {
			continue
		}

		for _, s05 := range cnt.S05 {
			endDate, err := p.formatDate(s05.Fh)

			if err != nil {
				continue
			}

			measureDateFile := p.getMeasureDateFromFilePath(dto.FilePath)
			measureDate, err := p.formatDate(measureDateFile)

			if err != nil {
				return []gross_measures.MeasureCloseWrite{}, err
			}

			distributorCdos := p.GetDistributorFromFilePath(dto.FilePath)

			measure := gross_measures.MeasureCloseWrite{
				EndDate:           endDate,
				ReadingDate:       measureDate.UTC(),
				Type:              measures.Absolute,
				ReadingType:       measures.DailyClosure,
				MeterSerialNumber: cnt.Id,
				ConcentratorID:    report.CncS05.Id,
				File:              dto.FilePath,
				DistributorCDOS:   distributorCdos,
				Origin:            measures.STG,
				Contract:          s05.Ctr,
				Periods: []gross_measures.MeasureClosePeriod{
					{
						Period: p.FormatPeriod(s05.Pt),
						AI:     s05.Value.AI * 1000,
						AE:     s05.Value.AE * 1000,
						R1:     s05.Value.R1 * 1000,
						R2:     s05.Value.R2 * 1000,
						R3:     s05.Value.R3 * 1000,
						R4:     s05.Value.R4 * 1000,
					},
				},
			}

			key := fmt.Sprintf("%s%v%v%s%s", measure.EndDate, measure.ReadingDate, measure.Contract, measure.MeterSerialNumber, measure.ConcentratorID)

			savedMeasure, ok := measuresParsed[key]

			if ok {
				measure.Periods = append(savedMeasure.Periods, measure.Periods...)
			}
			measuresParsed[key] = measure
		}

	}

	return utils.MapToSlice(measuresParsed), nil

}
