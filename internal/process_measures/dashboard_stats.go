package process_measures

import "context"

type ProcessMeasureDashboardStatsGlobal struct {
	DistributorID string       `json:"distributor_id" bson:"distributor_id"`
	Month         int          `json:"month" bson:"month"`
	Year          int          `json:"year" bson:"year"`
	Type          string       `json:"type" bson:"type"`
	Cups          string       `json:"cups" bson:"cups"`
	DailyStats    []DailyStats `json:"dailyStats" bson:"dailyStats"`
}
type DailyStats struct {
	Date                   int                  `json:"date" bson:"date"`
	ExpectedCurve          int                  `json:"expected_curve" bson:"expected_curve"`
	ExpectedDailyClosure   int                  `json:"expected_dailyClosure" bson:"expected_dailyClosure"`
	ExpectedMonthlyClosure int                  `json:"expected_monthlyClosure" bson:"expected_monthlyClosure"`
	LoadCurve              DashboardSingleStats `json:"load_curve" bson:"load_curve"`
	DailyClosure           DashboardSingleStats `json:"daily_closure" bson:"daily_closure"`
	MonthlyClosure         DashboardSingleStats `json:"monthly_closure" bson:"monthly_closure"`
}
type DashboardSingleStats struct {
	Total      int `json:"total" bson:"total"`
	Complete   int `json:"complete" bson:"complete"`
	Incomplete int `json:"incomplete" bson:"incomplete"`
	Absent     int `json:"absent" bson:"absent"`
	Val        int `json:"val" bson:"val"`
	Inv        int `json:"inv" bson:"inv"`
	Sup        int `json:"sup" bson:"sup"`
}

type SearchDashboardStats struct {
	DistributorID string
	Month         int
	Year          int
	Type          string
	Offset        int
	Limit         int
}

type ProcessMeasureDashboardStatsRepository interface {
	GetStatisticsGlobal(ctx context.Context, dto SearchDashboardStats) ([]ProcessMeasureDashboardStatsGlobal, error)
	GetStatisticsCups(ctx context.Context, q SearchDashboardStats) ([]ProcessMeasureDashboardStatsGlobal, int, error)
}
