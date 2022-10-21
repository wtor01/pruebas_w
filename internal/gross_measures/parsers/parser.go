package parsers

import (
	"archive/zip"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type ParserCurveI interface {
	Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCurveWrite, error)
	IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool
}

type ParserCloseI interface {
	Handle(ctx context.Context, dto gross_measures.HandleFileDTO) ([]gross_measures.MeasureCloseWrite, error)
	IsFile(ctx context.Context, dto gross_measures.HandleFileDTO) bool
}

type Parser struct {
	storage  storage.StorageCreator
	location *time.Location
}

func NewParser(storage storage.StorageCreator, loc *time.Location) Parser {
	return Parser{
		storage:  storage,
		location: loc,
	}
}

func (p Parser) GetFileNameFromFilePath(filePath string) string {
	return filepath.Base(filePath)
}

func (p Parser) GetFileExtension(file string) string {
	return filepath.Ext(file)
}

func (p Parser) FormatPeriod(period string) measures.PeriodKey {
	return measures.PeriodKey(fmt.Sprintf("P%s", period))
}

func (p Parser) FormatDateFromFile(date string) (time.Time, error) {
	m1 := regexp.MustCompile(`\D+`)
	d := m1.ReplaceAllString(date, "")
	if len(d) < 14 {
		return time.Time{}, errors.New("invalid date")
	}
	milliseconds := d[14:]
	layout := "20060102150405"
	if len(milliseconds) > 0 {
		layout = "20060102150405" + "." + strings.Repeat("0", utf8.RuneCountInString(milliseconds))
		d = d[0:14] + "." + d[14:]
	}

	return p.FormatDate(layout, d)
}

func (p Parser) FormatDate(layout string, date string) (time.Time, error) {
	locSpa, _ := time.LoadLocation("Europe/Madrid")

	return time.ParseInLocation(layout, date, locSpa)
}

func (p Parser) FormatDateBySeason(layout, date string, season rune) (time.Time, error) {
	parsedDate, err := time.Parse(layout, date)

	if season == 'S' || season == 'V' {
		parsedDate = parsedDate.Add(-time.Hour * 2)
	} else if season == 'W' || season == 'I' {
		parsedDate = parsedDate.Add(-time.Hour * 1)
	} else {
		err = errors.New("invalid season, S/W or V/I is valid")
	}

	return parsedDate, err
}

func (p Parser) GetMeasureDateFromFilePath(filePath string, datePosition int64) (time.Time, error) {
	filename := p.GetFileNameFromFilePath(filePath)
	filenamePaths := strings.Split(filename, "_")

	return p.FormatDateFromFile(filenamePaths[datePosition])
}

func (p Parser) GetDistributorFromFilePath(filePath string) string {
	paths := strings.Split(filePath, "/")
	if len(paths) < 2 {
		return ""
	}
	return paths[1]
}

func (p Parser) GetZipTmpFilePath(filePath string) string {
	filename := p.GetFileNameFromFilePath(filePath)
	return filepath.Join("/tmp/", filename)

}

func (p Parser) GetTmpFilePath(filePath string) string {
	filename := p.GetFileNameFromFilePath(filePath)
	return filepath.Join("/tmp/", strings.TrimSuffix(filename, filepath.Ext(filename)))
}

func (p Parser) Unzip(fileZip, dstFilepath string) error {
	archive, err := zip.OpenReader(fileZip)
	if err != nil {
		panic(err)
	}

	defer func() {
		archive.Close()
	}()
	for _, f := range archive.File {
		filePath := filepath.Join(dstFilepath, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(dstFilepath)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
