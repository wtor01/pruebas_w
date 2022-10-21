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

type PrimeS04 struct {
	Parser
}

type ReportS04 struct {
	XMLName xml.Name `xml:"Report"`
	CncS04  CncS04   `xml:"Cnc"`
}

type CncS04 struct {
	XMLName xml.Name `xml:"Cnc"`
	Id      string   `xml:"Id,attr"`
	CntS04  []CntS04 `xml:"Cnt"`
}

type CntS04 struct {
	XMLName xml.Name `xml:"Cnt"`
	Id      string   `xml:"Id,attr"`
	S04     []S04    `xml:"S04"`
	ErrCat  string   `xml:"ErrCat,attr"`
	ErrCode string   `xml:"ErrCode,attr"`
}

type S04 struct {
	XMLName xml.Name   `xml:"S04"`
	Fhi     string     `xml:"Fhi,attr"`
	Fhf     string     `xml:"Fhf,attr"`
	Ctr     string     `xml:"Ctr,attr"`
	Pt      string     `xml:"Pt,attr"`
	Mx      float64    `xml:"Mx,attr"`
	Fx      string     `xml:"Fx,attr"`
	Values  []ValueS04 `xml:"Value"`
}

type ValueS04 struct {
	XMLName xml.Name `xml:"Value"`
	AIa     float64  `xml:"AIa,attr"`
	AIi     float64  `xml:"AIi,attr"`
	AEa     float64  `xml:"AEa,attr"`
	AEi     float64  `xml:"AEi,attr"`
	R1a     float64  `xml:"R1a,attr"`
	R1i     float64  `xml:"R1i,attr"`
	R2a     float64  `xml:"R2a,attr"`
	R2i     float64  `xml:"R2i,attr"`
	R3a     float64  `xml:"R3a,attr"`
	R3i     float64  `xml:"R3i,attr"`
	R4a     float64  `xml:"R4a,attr"`
	R4i     float64  `xml:"R4i,attr"`
}

func NewPrimeS04(parser Parser) PrimeS04 {
	return PrimeS04{
		parser,
	}
}

func (p PrimeS04) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

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

	return filenamePaths[2] == "S04"
}

func (p PrimeS04) formatDate(date string) (time.Time, error) {
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

func (p PrimeS04) getMeasureDateFromFilePath(filePath string) string {
	filename := p.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return ""
	}

	return filenamePaths[4]
}

func (p PrimeS04) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error) {
	storage, err := p.storage(ctx)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	defer storage.Close()
	content, err := storage.ReadAll(ctx, dto.FilePath)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}
	var report ReportS04
	err = xml.Unmarshal(content, &report)
	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	measureDateFile := p.getMeasureDateFromFilePath(dto.FilePath)
	measureDate, err := p.formatDate(measureDateFile)

	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	distributorCdos := p.GetDistributorFromFilePath(dto.FilePath)

	measuresParsed := make(map[string]gross_measures.MeasureCloseWrite, 0)

	for _, cnt := range report.CncS04.CntS04 {
		if cnt.ErrCode != "" {
			continue
		}

		for _, s04 := range cnt.S04 {

			startDate, err := p.formatDate(s04.Fhi)

			if err != nil {
				continue
			}

			endDate, err := p.formatDate(s04.Fhf)

			if err != nil {
				continue
			}

			fxDate, err := p.formatDate(s04.Fx)

			if err != nil && s04.Mx > 0 {
				continue
			}

			if err != nil {
				fxDate = startDate
			}

			for i, values := range s04.Values {
				measure := gross_measures.MeasureCloseWrite{}
				if i == 0 {
					measure = gross_measures.MeasureCloseWrite{
						StartDate:         startDate,
						EndDate:           endDate,
						ReadingDate:       measureDate.UTC(),
						Type:              measures.Absolute,
						Contract:          s04.Ctr,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: cnt.Id,
						ConcentratorID:    report.CncS04.Id,
						File:              dto.FilePath,
						DistributorCDOS:   distributorCdos,
						Origin:            measures.STG,
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: p.FormatPeriod(s04.Pt),
								AI:     values.AIa * 1000,
								AE:     values.AEa * 1000,
								R1:     values.R1a * 1000,
								R2:     values.R2a * 1000,
								R3:     values.R3a * 1000,
								R4:     values.R4a * 1000,
								MX:     s04.Mx,
								FX:     fxDate,
							},
						},
					}
				} else {
					measure = gross_measures.MeasureCloseWrite{
						StartDate:         startDate,
						EndDate:           endDate,
						ReadingDate:       measureDate.UTC(),
						Type:              measures.Incremental,
						Contract:          s04.Ctr,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: cnt.Id,
						ConcentratorID:    report.CncS04.Id,
						File:              dto.FilePath,
						DistributorCDOS:   distributorCdos,
						Origin:            measures.STG,
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: p.FormatPeriod(s04.Pt),
								AI:     values.AIi * 1000,
								AE:     values.AEi * 1000,
								R1:     values.R1i * 1000,
								R2:     values.R2i * 1000,
								R3:     values.R3i * 1000,
								R4:     values.R4i * 1000,
								MX:     s04.Mx,
								FX:     fxDate,
							},
						},
					}
				}
				key := fmt.Sprintf("%s%v%v%s%s%s", measure.Type, measure.EndDate, measure.StartDate, measure.Contract, measure.MeterSerialNumber, measure.ConcentratorID)

				savedMeasure, ok := measuresParsed[key]

				if ok {
					measure.Periods = append(savedMeasure.Periods, measure.Periods...)
				}
				measuresParsed[key] = measure
			}

		}
	}
	return utils.MapToSlice(measuresParsed), nil
}
