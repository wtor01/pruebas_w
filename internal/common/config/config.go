package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProjectID                   string `envconfig:"PROJECT_ID" required:"true"`
	TopicMeasures               string `envconfig:"TOPIC_MEASURES" required:"false"`
	TopicBillingMeasures        string `envconfig:"TOPIC_BILLING_MEASURES" required:"false"`
	SubscriptionMeasures        string `envconfig:"SUBSCRIPTION_MEASURES" required:"false"`
	SubscriptionBillingMeasures string `envconfig:"SUBSCRIPTION_BILLING_MEASURES" required:"false"`
	StorageMeasuresFail         string `envconfig:"STORAGE_MEASURES_FAIL" required:"false"`
	StorageMeasuresSuccess      string `envconfig:"STORAGE_MEASURES_SUCCESS" required:"false"`

	TopicAggregations        string `envconfig:"TOPIC_AGGREGATIONS" required:"false"`
	SubscriptionAggregations string `envconfig:"SUBSCRIPTION_AGGREGATIONS" required:"false"`

	TopicCalendarPeriods        string `envconfig:"TOPIC_CALENDAR_PERIODS" required:"false"`
	SubscriptionCalendarPeriods string `envconfig:"SUBSCRIPTION_CALENDAR_PERIODS" required:"false"`

	Port            uint          `envconfig:"PORT" default:"8080"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`

	MeasureDB                string `envconfig:"MEASURE_DATABASE"`
	MeasureDBHost            string `envconfig:"MEASURE_DATABASE_HOST"`
	MeasureCollection        string `envconfig:"MEASURE_COLLECTION"`
	MeasureDBUser            string `envconfig:"MEASURE_DATABASE_USER"`
	MeasureDBPass            string `envconfig:"MEASURE_DATABASE_PASSWORD"`
	CollectionLoadCurve      string `envconfig:"COLLECTION_LOAD_CURVE"`
	CollectionDailyLoadCurve string `envconfig:"COLLECTION_DAILY_LOAD_CURVE"`
	CollectionMonthlyClosure string `envconfig:"COLLECTION_MONTHLY_CLOSURE"`
	CollectionDailyClosure   string `envconfig:"COLLECTION_DAILY_CLOSURE"`
	ProcessedMeasuresSub     string `envconfig:"PROCESSED_MEASURE_SUB"`
	ReProcessedMeasuresSub   string `envconfig:"RE_PROCESSED_MEASURE_SUB"`
	ProcessedMeasureTopic    string `envconfig:"PROCESSED_MEASURE_TOPIC"`

	SmarkiaHost  string `envconfig:"SMARKIA_HOST"`
	SmarkiaToken string `envconfig:"SMARKIA_TOKEN"`
	SmarkiaSub   string `envconfig:"SMARKIA_SUB"`
	SmarkiaTopic string `envconfig:"SMARKIA_TOPIC"`

	DbHost        string `envconfig:"POSTGRES_HOST"`
	DbPass        string `envconfig:"POSTGRES_PASSWORD"`
	DbUser        string `envconfig:"POSTGRES_USER"`
	DbName        string `envconfig:"POSTGRES_DB"`
	DbPort        string `envconfig:"POSTGRES_PORT"`
	RedisPort     string `envconfig:"REDIS_PORT"`
	RedisHost     string `envconfig:"REDIS_HOST"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisEnabled  bool   `envconfig:"REDIS_ENABLED"`

	// process measures
	Location string `envconfig:"LOCATION" default:"europe-west1"`
	TimeZone string `envconfig:"TIME_ZONE" default:"Europe/Madrid"`

	LocalLocation *time.Location
}

func LoadConfig(logger *log.Logger) (Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Println(err)
	}
	envPrefix := os.Getenv("ENV_PREFIX")
	var cfg Config
	err = envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return Config{}, err
	}

	loc, err := time.LoadLocation(cfg.TimeZone)

	if err != nil {
		logger.Fatalln(err)
	}

	cfg.LocalLocation = loc
	return cfg, err
}
