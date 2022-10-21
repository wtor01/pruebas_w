package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"time"
)

type GetDashboardQuery struct {
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

type MeasureServicePointDashboard struct {
	Status measures.Status  `json:"status"`
	P0     *measures.Values `json:"P0"`
	P1     *measures.Values `json:"P1"`
	P2     *measures.Values `json:"P2"`
	P3     *measures.Values `json:"P3"`
	P4     *measures.Values `json:"P4"`
	P5     *measures.Values `json:"P5"`
	P6     *measures.Values `json:"P6"`
}

func (m *MeasureServicePointDashboard) GetPeriodMeasure(period measures.PeriodKey) *measures.Values {

	switch period {
	case measures.P0:
		return m.P0
	case measures.P1:
		return m.P1
	case measures.P2:
		return m.P2
	case measures.P3:
		return m.P3
	case measures.P4:
		return m.P4
	case measures.P5:
		return m.P5
	case measures.P6:
		return m.P6
	}

	return nil
}

func (m *MeasureServicePointDashboard) SetPeriodMeasure(period measures.PeriodKey, measure *measures.Values) {

	switch period {
	case measures.P0:
		m.P0 = measure
	case measures.P1:
		m.P1 = measure
	case measures.P2:
		m.P2 = measure
	case measures.P3:
		m.P3 = measure
	case measures.P4:
		m.P4 = measure
	case measures.P5:
		m.P5 = measure
	case measures.P6:
		m.P6 = measure
	}

}

func (m *MeasureServicePointDashboard) SetStatus(status measures.Status) {
	if m.Status.Compare(measures.StatusValue, status) {
		m.Status = status
	}
}

// SetStatusProcessCurve If compare returns false, we set the status of the new value to the previous one.
// Compare the previous status with new status process curve value and return true or false
func (m *MeasureServicePointDashboard) SetStatusProcessCurve(status measures.Status) {
	if !m.Status.Compare(measures.StatusProcessCurveValue, status) {
		m.Status = status
	}
}

type MeasureServicePointDashboardMonthly struct {
	InitDate string                  `json:"init_date"`
	EndDate  string                  `json:"end_date"`
	ID       string                  `json:"id"`
	Status   measures.Status         `json:"status"`
	P0       *measures.ValuesMonthly `json:"P0"`
	P1       *measures.ValuesMonthly `json:"P1"`
	P2       *measures.ValuesMonthly `json:"P2"`
	P3       *measures.ValuesMonthly `json:"P3"`
	P4       *measures.ValuesMonthly `json:"P4"`
	P5       *measures.ValuesMonthly `json:"P5"`
	P6       *measures.ValuesMonthly `json:"P6"`
}

func (m *MeasureServicePointDashboardMonthly) SetStatus(status measures.Status) {
	if m.Status == "" {
		m.Status = status
	}
	if m.Status.Compare(measures.StatusValue, status) {
		m.Status = status
	}
}

func (m *MeasureServicePointDashboardMonthly) GetPeriodMeasure(period measures.PeriodKey) *measures.ValuesMonthly {

	switch period {
	case measures.P0:
		return m.P0
	case measures.P1:
		return m.P1
	case measures.P2:
		return m.P2
	case measures.P3:
		return m.P3
	case measures.P4:
		return m.P4
	case measures.P5:
		return m.P5
	case measures.P6:
		return m.P6
	}

	return nil
}

func (m *MeasureServicePointDashboardMonthly) SetPeriodMeasure(period measures.PeriodKey, measure *measures.ValuesMonthly) {

	switch period {
	case measures.P0:
		m.P0 = measure
	case measures.P1:
		m.P1 = measure
	case measures.P2:
		m.P2 = measure
	case measures.P3:
		m.P3 = measure
	case measures.P4:
		m.P4 = measure
	case measures.P5:
		m.P5 = measure
	case measures.P6:
		m.P6 = measure
	}

}

type ServicePointDashboardWithType struct {
	Type          measures.MeterType
	ServicePoints []ServicePointDashboard
}
type ServicePointDashboard struct {
	Date            string                              `json:"date"`
	MagnitudeEnergy measures.Magnitude                  `json:"magnitude_energy"`
	Magnitudes      []measures.Magnitude                `json:"magnitudes"`
	Periods         []measures.PeriodKey                `json:"periods"`
	MonthlyClosure  MeasureServicePointDashboardMonthly `json:"monthly_closure"`
	DailyClosure    MeasureServicePointDashboard        `json:"daily_closure"`
	Curve           MeasureServicePointDashboard        `json:"curve"`
}

func (s *ServicePointDashboard) SetCurves(measure MeasureServicePointDashboard) {
	s.Curve = measure
}

func (s *ServicePointDashboard) SetDaily(measure MeasureServicePointDashboard) {
	s.DailyClosure = measure
}

func (s *ServicePointDashboard) SetMonthly(measure MeasureServicePointDashboardMonthly) {
	s.MonthlyClosure = measure
}

type ServicePointDashboardCurve struct {
	Date   time.Time       `json:"date"`
	Status measures.Status `json:"status"`
	Values struct {
		measures.Values
	} `json:"values"`
}

type ListCupsQuery struct {
	Cups      []string
	StartDate time.Time
	EndDate   time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ProcessMeasureDashboardRepository
type ProcessMeasureDashboardRepository interface {
	GetDashboard(ctx context.Context, query GetDashboardQuery) ([]measures.DashboardMeasureI, error)
	GetCupsMeasures(ctx context.Context, query ListCupsQuery) (map[string]*measures.DashboardCupsReading, error)
}
