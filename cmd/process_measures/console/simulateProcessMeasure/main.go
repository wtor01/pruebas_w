package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	validation_services "bitbucket.org/sercide/data-ingestion/internal/validations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

func main() {

	logger := log.Default()

	cnf, _ := config.LoadConfig(logger)

	distributorId := "5140cda8-1daa-4f52-b85d-bf29d45ca62a"
	cups := "ES0130000001348006SX"
	date := time.Date(2022, 10, 1, 0, 0, 0, 0, cnf.LocalLocation)
	simulateProcessMeasure(distributorId, cups, date, measures.DailyClosure)

}

func simulateProcessMeasure(distributorId string, cups string, date time.Time, processType measures.ReadingType) {

	logger := log.Default()

	cnf, err := config.LoadConfig(logger)

	ctx := context.Background()

	mongoClient, err := mongo_repos.New(ctx, cnf)

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var rd *redis.Client

	if cnf.RedisEnabled {
		rd = redis_repos.New(cnf)
		defer rd.Close()
	}

	db := postgres.New(cnf)

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	schedulerRepository := postgres.NewProcessMeasureSchedulerPostgres(db)
	adminRepository := postgres.NewValidationMeasurePostgresRepository(db, redis_repos.NewDataCacheRedis(rd), cnf.ShutdownTimeout)
	calendarRepository := postgres.NewCalendarPostgres(db, redis_repos.NewDataCacheRedis(rd))
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(rd))
	seasonRepository := postgres.NewSeasonPostgres(db, redis_repos.NewDataCacheRedis(rd))

	processRepository := mongo_repos.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)
	validationRepository := mongo_repos.NewValidationMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	grossRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	inventoryMeasuresRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	repositoryFestiveDays := redis_repos.NewProcessMeasureFestiveDays(rd)

	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))

	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	srv := services.NewServicesPubsub(
		publisher,
		cnf.ProcessedMeasureTopic,
		cnf.TopicBillingMeasures,
		schedulerRepository,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryMeasuresRepository), redis_repos.NewDataCacheRedis(rd)),
		inventoryMeasuresRepository,
		grossRepository,
		processRepository,
		repositoryFestiveDays,
		seasonRepository,
		cnf.LocalLocation,
		internal_clients.NewValidation(validation_services.NewServices(adminRepository, processRepository)),
		calendarRepository,
		masterTablesClient,
		validationRepository,
	)

	meterConfig, err := inventoryMeasuresRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		Distributor: distributorId,
		CUPS:        cups,
		Time:        date,
	})

	msg := getMessageAttributes(processType, meterConfig, date)

	switch processType {
	case measures.Curve:
		srv.ProcessCurve.Handle(ctx, msg.Payload)
	case measures.DailyClosure:
		srv.ProcessDailyClosure.Handle(ctx, msg.Payload)
	case measures.BillingClosure:
		srv.ProcessMonthlyClosure.Handle(ctx, msg.Payload)
	}
}

func getMessageAttributes(processType measures.ReadingType, meterConfig measures.MeterConfig, date time.Time) measures.ProcessMeasureEvent {
	var msg measures.ProcessMeasureEvent

	switch processType {
	case measures.Curve:
		{
			if utils.InSlice(meterConfig.CurveType, []measures.RegisterType{measures.Hourly, measures.Both}) {
				msg = process_measures.NewProcessCurveEvent(date, meterConfig, measures.HourlyMeasureCurveReadingType)
			}

			if utils.InSlice(meterConfig.CurveType, []measures.RegisterType{measures.QuarterHour, measures.Both}) {
				msg = process_measures.NewProcessCurveEvent(date, meterConfig, measures.QuarterMeasureCurveReadingType)
			}

		}
	case measures.BillingClosure:
		{
			msg = process_measures.NewProcessBillingClosureEvent(date, meterConfig)
		}
	case measures.DailyClosure:
		{
			msg = process_measures.NewProcessDailyClosureEvent(date, meterConfig)
		}
	}

	return msg
}
