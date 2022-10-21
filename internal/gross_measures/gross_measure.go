package gross_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidDistributorId = errors.New("invalid distributor id")
	ErrInvalidMeasureData   = errors.New("invalid measure data")
	ErrInvalidMeasureDate   = errors.New("invalid measure date")
)

type GrossMeasureBase interface {
	GenerateID()
	GetID() string
	SetDistributorID(distributorId string)
	GetMeterSerialNumber() string
	GetOriginFile() string
	GetReadingType() measures.ReadingType
	GetEndDate() time.Time
	GetDistributorCDOS() string
	GetDistributorID() string
	ToValidatable() []validations.MeasureValidatable
	SetStatus(status measures.Status)
	SetStatusMeasure(validation validations.ValidatorBase)
}

type HandleFileDTO struct {
	FilePath string
}

type MeasureService interface {
	HandleFile(ctx context.Context, dto HandleFileDTO) error
}

type GetMeasuresQuery struct {
	DistributorId string
	MeterId       string
	StartDate     time.Time
	EndDate       time.Time
	ReadingType   string
	Origin        string
}

type QueryListForProcessCurve struct {
	SerialNumber string
	Date         time.Time
	CurveType    measures.MeasureCurveReadingType
}

type QueryListForProcessCurveGenerationDate struct {
	ReadingType   measures.ReadingType
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
	Limit         int
	Offset        int
}

type QueryListForProcessClose struct {
	SerialNumber string
	Date         time.Time
	ReadingType  measures.ReadingType
}

type QueryListMeasure struct {
	SerialNumber string
	StartDate    time.Time
	EndDate      time.Time
	ReadingType  measures.ReadingType
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=GrossMeasureRepository
type GrossMeasureRepository interface {
	SaveAllMeasuresCurve(ctx context.Context, measures []MeasureCurveWrite) error
	SaveAllMeasuresClose(ctx context.Context, measures []MeasureCloseWrite) error
	ListDailyCurveMeasures(ctx context.Context, query QueryListForProcessCurve) ([]MeasureCurveWrite, error)
	ListDailyCloseMeasures(ctx context.Context, query QueryListForProcessClose) ([]MeasureCloseWrite, error)
	ListCurveMeasures(ctx context.Context, query QueryListMeasure) ([]MeasureCurveWrite, error)
	ListCloseMeasures(ctx context.Context, query QueryListMeasure) ([]MeasureCloseWrite, error)
	ListGrossMeasuresFromGenerationDate(ctx context.Context, query QueryListForProcessCurveGenerationDate) ([]MeasureCurveMeterSerialNumber, error)
	CountGrossMeasuresFromGenerationDate(ctx context.Context, query QueryListForProcessCurveGenerationDate) (int, error)
}
