package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"encoding/xml"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type PrimeS02 struct {
	Parser
}

type Report struct {
	XMLName xml.Name `xml:"Report"`
	Cnc     Cnc      `xml:"Cnc"`
}

type Cnc struct {
	XMLName xml.Name `xml:"Cnc"`
	Id      string   `xml:"Id,attr"`
	Cnt     []Cnt    `xml:"Cnt"`
}

type Cnt struct {
	XMLName xml.Name `xml:"Cnt"`
	Id      string   `xml:"Id,attr"`
	S02     []S02    `xml:"S02"`
	ErrCat  string   `xml:"ErrCat,attr"`
	ErrCode string   `xml:"ErrCode,attr"`
	Magn    int64    `xml:"Magn,attr"`
}

type S02 struct {
	XMLName xml.Name `xml:"S02"`

	Fh string  `xml:"Fh,attr"`
	AI float64 `xml:"AI,attr"`
	AE float64 `xml:"AE,attr"`
	R1 float64 `xml:"R1,attr"`
	R2 float64 `xml:"R2,attr"`
	R3 float64 `xml:"R3,attr"`
	R4 float64 `xml:"R4,attr"`
	Bc string  `xml:"Bc,attr"`
}

func (p PrimeS02) IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool {

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

	return filenamePaths[2] == "S02"
}

func NewPrimeS02(parser Parser) PrimeS02 {
	return PrimeS02{
		parser,
	}
}

func (p PrimeS02) formatDate(date string) (time.Time, error) {
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

func (p PrimeS02) getMeasureDateFromFilePath(filePath string) string {
	filename := p.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")
	if len(filenamePaths) != 5 {
		return ""
	}

	return filenamePaths[4]
}

func (p PrimeS02) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCurveWrite, error) {
	storage, err := p.storage(ctx)
	if err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}

	defer storage.Close()
	content, err := storage.ReadAll(ctx, dto.FilePath)
	if err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}
	var report Report
	err = xml.Unmarshal(content, &report)
	if err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}
	measureList := make([]gross_measures.MeasureCurveWrite, 0)

	for _, cnt := range report.Cnc.Cnt {
		if cnt.ErrCode != "" {
			continue
		}

		for _, s02 := range cnt.S02 {
			endDate, err := p.formatDate(s02.Fh)

			if err != nil {
				continue
			}

			measureDateFile := p.getMeasureDateFromFilePath(dto.FilePath)
			measureDate, err := p.formatDate(measureDateFile)

			if err != nil {
				return []gross_measures.MeasureCurveWrite{}, err
			}

			distributorCdos := p.GetDistributorFromFilePath(dto.FilePath)

			if distributorCdos == "" {
				return []gross_measures.MeasureCurveWrite{}, errors.New("invalid distributor id")
			}

			if cnt.Magn == 0 {
				cnt.Magn = 1
			}

			measureList = append(measureList, gross_measures.MeasureCurveWrite{
				EndDate:           endDate,
				ReadingDate:       measureDate.UTC(),
				Type:              measures.Incremental,
				ReadingType:       measures.Curve,
				CurveType:         measures.HourlyMeasureCurveReadingType,
				MeterSerialNumber: cnt.Id,
				ConcentratorID:    report.Cnc.Id,
				File:              dto.FilePath,
				DistributorCDOS:   distributorCdos,
				Origin:            measures.STG,
				Qualifier:         s02.Bc,
				AI:                s02.AI * float64(cnt.Magn),
				AE:                s02.AE * float64(cnt.Magn),
				R1:                s02.R1 * float64(cnt.Magn),
				R2:                s02.R2 * float64(cnt.Magn),
				R3:                s02.R3 * float64(cnt.Magn),
				R4:                s02.R4 * float64(cnt.Magn),
			})
		}
	}
	return measureList, nil
}
