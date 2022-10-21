package gross_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"time"
)

type CalendarStatus struct {
	Date   string          `json:"date"`
	Status measures.Status `json:"status"`
}

func (c *CalendarStatus) SetStatus(status measures.Status) {
	if !c.Status.Compare(measures.StatusValue, status) {
		c.Status = status
	}
}

type ListDailyClose struct {
	EndDate string          `json:"end_date"`
	Status  measures.Status `json:"status"`
	Origin  string          `json:"origin"`
	P0      measures.Values `json:"p0"`
	P1      measures.Values `json:"p1"`
	P2      measures.Values `json:"p2"`
	P3      measures.Values `json:"p3"`
	P4      measures.Values `json:"p4"`
	P5      measures.Values `json:"p5"`
	P6      measures.Values `json:"p6"`
}

func (c *ListDailyClose) SetPeriodValues(period measures.PeriodKey, values measures.Values) {
	switch period {
	case measures.P0:
		c.P0 = values
	case measures.P1:
		c.P1 = values
	case measures.P2:
		c.P2 = values
	case measures.P3:
		c.P3 = values
	case measures.P4:
		c.P4 = values
	case measures.P5:
		c.P5 = values
	case measures.P6:
		c.P6 = values
	}
}

type ListMonthlyClose struct {
	InitDate string                 `json:"init_date"`
	EndDate  string                 `json:"end_date"`
	Status   measures.Status        `json:"status"`
	Origin   string                 `json:"origin"`
	P0       measures.ValuesMonthly `json:"p0"`
	P1       measures.ValuesMonthly `json:"p1"`
	P2       measures.ValuesMonthly `json:"p2"`
	P3       measures.ValuesMonthly `json:"p3"`
	P4       measures.ValuesMonthly `json:"p4"`
	P5       measures.ValuesMonthly `json:"p5"`
	P6       measures.ValuesMonthly `json:"p6"`
}

func (c *ListMonthlyClose) SetPeriodValues(period measures.PeriodKey, values measures.Values) {
	periodValue := c.GetPeriod(period)
	monthlyValue := measures.ValuesMonthly{
		Values: values,
	}

	if periodValue != nil {
		monthlyValue.AIi = periodValue.AIi
		monthlyValue.AEi = periodValue.AEi
		monthlyValue.R1i = periodValue.R1i
		monthlyValue.R2i = periodValue.R2i
		monthlyValue.R3i = periodValue.R3i
		monthlyValue.R4i = periodValue.R4i
	}

	switch period {
	case measures.P0:
		c.P0 = monthlyValue
	case measures.P1:
		c.P1 = monthlyValue
	case measures.P2:
		c.P2 = monthlyValue
	case measures.P3:
		c.P3 = monthlyValue
	case measures.P4:
		c.P4 = monthlyValue
	case measures.P5:
		c.P5 = monthlyValue
	case measures.P6:
		c.P6 = monthlyValue
	}
}

func (c *ListMonthlyClose) SetPeriodMonthlyValues(period measures.PeriodKey, values measures.Values) {
	periodValue := c.GetPeriod(period)

	monthlyValue := measures.ValuesMonthly{
		AIi: values.AI,
		AEi: values.AE,
		R1i: values.R1,
		R2i: values.R2,
		R3i: values.R3,
		R4i: values.R4,
	}

	if periodValue != nil {
		monthlyValue.Values = periodValue.Values
	}

	switch period {
	case measures.P0:
		c.P0 = monthlyValue
	case measures.P1:
		c.P1 = monthlyValue
	case measures.P2:
		c.P2 = monthlyValue
	case measures.P3:
		c.P3 = monthlyValue
	case measures.P4:
		c.P4 = monthlyValue
	case measures.P5:
		c.P5 = monthlyValue
	case measures.P6:
		c.P6 = monthlyValue
	}
}

func (c *ListMonthlyClose) GetPeriod(period measures.PeriodKey) *measures.ValuesMonthly {
	switch period {
	case measures.P0:
		return &c.P0
	case measures.P1:
		return &c.P1
	case measures.P2:
		return &c.P2
	case measures.P3:
		return &c.P3
	case measures.P4:
		return &c.P4
	case measures.P5:
		return &c.P5
	case measures.P6:
		return &c.P6

	}

	return nil
}

type DashboardMeasureSupplyPoint struct {
	Cups                   string               `json:"cups"`
	SerialNumber           string               `json:"serial_number"`
	Magnitudes             []measures.Magnitude `json:"magnitudes"`
	MagnitudeEnergy        measures.Magnitude   `json:"magnitude_energy"`
	MeterType              measures.MeterType   `json:"meter_type"`
	Periods                []measures.PeriodKey `json:"periods"`
	CalendarDailyClosure   []CalendarStatus     `json:"calendar_daily_closure"`
	CalendarMonthlyClosure []CalendarStatus     `json:"calendar_monthly_closure"`
	CalendarCurve          []CalendarStatus     `json:"calendar_curve"`
	ListDailyClosure       []ListDailyClose     `json:"list_daily_closure"`
	ListMonthlyClosure     []ListMonthlyClose   `json:"list_monthly_closure"`
}

func (d *DashboardMeasureSupplyPoint) SetCalendar(readingType measures.ReadingType, calendar []CalendarStatus) {
	switch readingType {
	case measures.DailyClosure:
		d.CalendarDailyClosure = calendar
	case measures.Curve:
		d.CalendarCurve = calendar
	case measures.BillingClosure:
		d.CalendarMonthlyClosure = calendar
	}
}

func (d *DashboardMeasureSupplyPoint) SetListDailyClosure(dailyClosures []ListDailyClose) {
	d.ListDailyClosure = dailyClosures
}

func (d *DashboardMeasureSupplyPoint) SetListMonthlyClosure(monthlyClosures []ListMonthlyClose) {
	d.ListMonthlyClosure = monthlyClosures
}

func NewDashboardMeasureSupplyPoint(config measures.MeterConfig, periods []measures.PeriodKey) DashboardMeasureSupplyPoint {
	return DashboardMeasureSupplyPoint{
		Cups:            config.Cups(),
		SerialNumber:    config.SerialNumber(),
		MagnitudeEnergy: config.EnergyMagnitude(),
		Magnitudes:      config.GetMagnitudesActive(),
		MeterType:       config.MeterType(),
		Periods:         periods,
	}
}

type DashboardSupplyPointCurve struct {
	Hour   string          `json:"hour"`
	Status measures.Status `json:"status"`
	File   string          `json:"file"`
	Values measures.Values `json:"values"`
}

type GetDashboardQuery struct {
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

type GetDashboardBySupplyPointQuery struct {
	DistributorId string
	SerialNumber  string
	StartDate     time.Time
	EndDate       time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=DashboardRepository
type DashboardRepository interface {
	GetDashboard(ctx context.Context, query GetDashboardQuery) ([]measures.DashboardMeasureI, error)
}
