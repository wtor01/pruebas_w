package gross_measures

import "context"

type DashboardSerialNumber struct {
	DistributorId    string                 `bson:"distributor_id" json:"distributor_id"`
	Month            int                    `bson:"month" json:"month"`
	Year             int                    `bson:"year" json:"year"`
	Cups             string                 `bson:"cups" json:"cups"`
	Type             string                 `bson:"type" json:"type"`
	SerialNumber     string                 `bson:"serial_number" json:"serial_number"`
	ServiceType      string                 `bson:"service_type" json:"service_type"`
	ServicePointType string                 `bson:"service_point_type" json:"service_point_type"`
	DailyStats       []DashboardSingleStats `bson:"dailyStats" json:"dailyStats"`
}
type GrossMeasuresDashboardStatsGlobal struct {
	DistributorID string                 `json:"distributor_id" bson:"distributor_id"`
	Month         int                    `json:"month" bson:"month"`
	Year          int                    `json:"year" bson:"year"`
	Type          string                 `json:"type" bson:"type"`
	DailyStats    []DashboardSingleStats `json:"dailyStats" bson:"dailyStats"`
}

type DashboardSingleStats struct {
	Day                    int `json:"date" bson:"date"`
	HourlyCurve            int `json:"hourly_curve" bson:"hourly_curve"`
	QuarterlyCurve         int `json:"quarterly_curve" bson:"quarterly_curve"`
	DailyClosure           int `json:"daily_closure" bson:"daily_closure"`
	MonthlyClosure         int `json:"monthly_closure" bson:"monthly_closure"`
	ExpectedHourlyCurve    int `json:"expected_hourly_curve" bson:"expected_hourly_curve"`
	ExpectedQuarterlyCurve int `json:"expected_quarterly_curve" bson:"expected_quarterly_curve"`
	ExpectedDailyClosure   int `json:"expected_daily_closure" bson:"expected_daily_closure"`
	ExpectedMonthlyClosure int `json:"expected_monthly_closure" bson:"expected_monthly_closure"`
}

type SearchDashboardSerialNumber struct {
	DistributorId string
	Month         int
	Year          int
	Type          string
	Ghost         bool
	Offset        int
	Limit         int
}

type SearchDashboardStats struct {
	DistributorID string
	Month         int
	Year          int
	Type          string
}

type ListGrossMeasuresStatisticsSerialNumberResult struct {
	Data  []DashboardSerialNumber `bson:"data"`
	Count int                     `bson:"total"`
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=GrossMeasuresDashboardStatsRepository
type GrossMeasuresDashboardStatsRepository interface {
	ListGrossMeasuresStatisticsSerialNumber(ctx context.Context, q SearchDashboardSerialNumber) (ListGrossMeasuresStatisticsSerialNumberResult, error)
	GetStatisticsGlobal(ctx context.Context, dto SearchDashboardStats) ([]GrossMeasuresDashboardStatsGlobal, error)
}
