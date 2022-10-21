package postgres

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

func Migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	db.Exec("CREATE TYPE service_point_point_type AS ENUM ('1', '2', '3', '4', '5');")

	db.Exec("CREATE TYPE service_point_service_type AS ENUM ('G-D', 'D-D', 'D-C');")

	db.Exec("CREATE TYPE service_point_point_tension_level AS ENUM ('AT2', 'MT', 'BT');")

	db.Exec("CREATE TYPE service_point_measure_tension_level AS ENUM ('AT2', 'MT', 'BT');")

	db.Exec("CREATE TYPE measure_points_type AS ENUM ('P', 'R', 'C', 'T');")

	db.Exec("CREATE TYPE meters_technology AS ENUM ('E', 'M', 'R', 'T');")

	db.Exec("CREATE TYPE meters_type AS ENUM ('TLG', 'TLM', 'OTHER');")

	db.Exec("CREATE TYPE meter_configs_curve_type AS ENUM ('HOURLY', 'QUARTER', 'BOTH', 'NONE');")

	db.Exec("CREATE TYPE meter_configs_reading_type AS ENUM ('INC', 'ABS');")

	db.Exec("CREATE TYPE meter_configs_contract_number AS ENUM ('1', '2', '3');")

	db.Exec("CREATE TYPE validation_rules_reading_type AS ENUM ('INC', 'ABS', 'INC_CLO', 'ABS_CLO');")

	db.Exec("CREATE TYPE validation_rules_type AS ENUM ('INM', 'PROC', 'COHE');")

	db.Exec("CREATE TYPE validation_rules_action AS ENUM ('ALERT', 'SUPERV', 'INV', 'NONE');")

	db.Exec("CREATE TYPE meter_configs_tlg_type AS ENUM ('TLG_OP_CURVE', 'NO_TLG', 'TLG_OP_NOCURVE', 'TLG_NOOP');")

	err := db.AutoMigrate(
		&BillingMeasuresScheduler{},
		&Distributor{},
		&MeasurePoint{},
		&Meter{},
		&User{},
		&Role{},
		&MeterConfig{},
		&ServicePoint{},
		&UserDistributorRole{},
		&ValidationRule{},
		&ValidationRuleConfig{},
		&HistoryValidationRuleConfig{},
		&ProcessMeasureScheduler{},
		&GeographicZones{},
		&Calendars{},
		&CalendarPeriods{},
		&Tariff{},
		&TariffCalendar{},
		&Seasons{},
		&DayTypes{},
		&FestiveDays{},
		&ContractualSituation{},
		&RecoverMeasures{},
		&AggregationsFeatures{},
		&AggregationsConfig{},
	)

	return err
}

func Rollback(db *gorm.DB) {
	types := []string{
		"service_point_point_type",
		"service_point_service_type",
		"service_point_point_tension_level",
		"service_point_measure_tension_level",
		"measure_points_type",
		"meters_technology",
		"meters_type",
		"meter_configs_curve_type",
		"meter_configs_reading_type",
		"meter_configs_contract_number",
		"validation_rules_reading_type",
		"validation_rules_type",
		"validation_rules_action",
		"process_measure_schedulers",
	}
	tables := []string{
		"distributors",
		"meters",
		"meter_configs",
		"measure_points",
		"roles",
		"service_points",
		"user_distributor_roles",
		"users",
		"validation_rule_configs",
		"validation_rules",
		"hist_validation_rules_configs",
		"geographic_zones",
		"calendars",
		"calendar_periods",
		"tariffs",
		"tariff_calendar",
	}

	for _, t := range tables {
		tx := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", t))
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}

	for _, t := range types {
		tx := db.Exec(fmt.Sprintf("DROP TYPE %s", t))
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}
}

func createServicePointConfigView(db *gorm.DB) error {
	tx := db.Exec(`

		CREATE MATERIALIZED VIEW IF NOT EXISTS service_point_config_view AS
			SELECT meter_configs.id                 AS meter_config_id,
				   meter_configs.reading_type       AS reading_type,
				   meter_configs.priority_contract,
				   meter_configs.calendar_id        AS calendar,
				   meter_configs.start_date,
				   meter_configs.end_date,
				   meter_configs.tlg_code,
				   meter_configs.tlg_type,
				   meter_configs.ai,
				   meter_configs.ae,
				   meter_configs.r1,
				   meter_configs.r2,
				   meter_configs.r3,
				   meter_configs.r4,
				   meters.id                        AS meter_id,
				   distributors.cdos                AS distributor_code,
				   distributors.id                  AS distributor_id,
				   service_points.cups              AS cups,
				   meters.serial_number             AS meter_serial_number,
				   meters.type                      AS meters_type,
				   service_points.service_type      AS service_type,
				   service_points.point_type        AS point_type,
				   measure_points.type              AS measure_point_type,
				   case
					   when (meters.type = 'OTHER') THEN 'NONE'
					   when (meters.type = 'TLG') THEN 'HOURLY'
					   else 'QUARTER' end           AS curve_type,
				   meters.technology AS technology,
				   contractual_situations.id        AS contractual_situation_id,
				   tariff_id,
				   service_points.id                AS service_point_id,
				   service_points.cups              AS service_point_cups,
				   tariffs.coef                     AS coefficient,
				   tariffs.calendar_id              AS calendar_code,
				   tariffs.geographic_id            AS geographic_id,
				   p1_demand,
				   p2_demand,
				   p3_demand,
				   p4_demand,
				   p5_demand,
				   p6_demand,
				   contractual_situations.init_date AS contractual_situation_start_date,
				   contractual_situations.end_date  AS contractual_situation_end_date
			FROM distributors
					 INNER JOIN service_points ON service_points.distributor_id = distributors.id
					 INNER JOIN contractual_situations ON contractual_situations.service_point_id = service_points.id
					 INNER JOIN tariffs ON contractual_situations.tariff_id = tariffs.id
					 INNER JOIN measure_points ON measure_points.service_point_id = service_points.id
					 INNER JOIN meter_configs ON meter_configs.measure_point_id = measure_points.id
					 INNER JOIN meters meters ON meter_configs.meter_id = meters.id;
	`)

	return tx.Error
}

func createRoles(db *gorm.DB) {

	me := &Role{
		ModelEntity: ModelEntity{
			ID:        uuid.FromStringOrNil("d03207f1-7bbd-4adb-bc28-85390131eee5"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: "admin",
	}

	user1 := &Role{
		ModelEntity: ModelEntity{
			ID:        uuid.FromStringOrNil("5ac2a9fe-119e-4311-90a5-015110184cd0"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: "operator",
	}

	db.Save(&me)
	db.Save(&user1)
}
