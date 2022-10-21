package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"os"
	"time"
)

func main() {
	//sendEventOnSaveBillingMeasure()
	//insertConsumProfile()
	processTgMeasureByCup()
}

func searchConsumProfile() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	repository := mongo.NewConsumProfileRepository(mongoClient, cnf.MeasureDB)

	result, err := repository.Search(context.TODO(), billing_measures.QueryConsumProfile{
		EndDate:   time.Date(2022, 01, 31, 00, 0, 0, 0, time.UTC),
		StartDate: time.Date(2022, 01, 01, 00, 0, 0, 0, time.UTC),
	})

	fmt.Println(result, err)
}

func insertConsumProfile() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	repository := mongo.NewConsumProfileRepository(mongoClient, cnf.MeasureDB)

	svc := services.NewInsertConsumProfile(storage.NewStorageGCP, repository)

	err = svc.Handler(context.Background(), services.InsertConsumProfileDTO{
		File: "eu.artifacts.med-n-base-fb6f.appspot.com/Datos_coeficientes_consumo.csv",
	})

	fmt.Println(err)
}

func processTgMeasureByCup() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}

	db := postgres.New(cnf)
	inventoryMeasuresRepository := mongo.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(redisClient))
	billingMeasuresRepository := mongo.NewBillingMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, "billing_measure")
	processRepository := mongo.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(redisClient))

	inventoryClient := inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryMeasuresRepository), redis_repos.NewDataCacheRedis(redisClient))

	repositoryFestiveDays := redis_repos.NewProcessMeasureFestiveDays(redisClient)
	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	svc := services.NewProcessMvhByCup(
		billingMeasuresRepository,
		processRepository,
		inventoryClient,
		mongo.NewConsumProfileRepository(mongoClient, cnf.MeasureDB),
		repositoryFestiveDays,
		cnf.LocalLocation,
		masterTablesClient,
		publisher,
		cnf.TopicBillingMeasures,
	)

	serviceName := "billing_measures"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)

	tracer := telemetry.SetTracer("billing_measures/" + serviceName + "/console")

	ctx, span := tracer.Start(ctx, "ProcessMvhByCup")

	defer span.End()

	date := time.Date(2022, 7, 31, 0, 0, 0, 0, cnf.LocalLocation).UTC()

	meterConfig, err := inventoryMeasuresRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        "ES0130000000300002JR",
		Time:        date,
		Distributor: "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
	})

	err = svc.Handler(ctx, measures.ProcessMeasurePayload{
		Date:        date,
		MeterConfig: meterConfig,
	})

	fmt.Println(err)
}

func sendEventOnSaveBillingMeasure() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	err = event.PublishEvent(ctx, cnf.TopicBillingMeasures, publisher, billing_measures.NewOnSaveBillingMeasureEvent(billing_measures.OnSaveBillingMeasurePayload{
		InitDate:         time.Time{},
		EndDate:          time.Time{},
		CUPS:             "ES0130000000343073EQ",
		BillingMeasureId: "e1d7b9a6121d233e16b10f7107923c1144a4be9aac7058c12c89a4904369a201",
	}))
	/*{
		"type": "BILLING_MEASURES/SELF_CONSUMPTION",
		"payload": {
		"init_date": "2022-04-30T22:00:00Z",
			"end_date": "2022-08-01T22:00:00Z",
			"CUPS": "ES0130000000314144YK",
			"billing_measure_id": "bd4805f6b68e147e1206bb238f4eb824f10e040db5bcbe73ac1242247a589dba"
	}
	}*/

	fmt.Println(err)
}

/*
CUPS IDS

5140cda8-1daa-4f52-b85d-bf29d45ca62a

TEST 1
ES0130000000343026CS 2d184517da7ae1b2ac3a152aed628730d8a28a7900610cd78586e5c663429ea0
ES0130000000343073EQ e1d7b9a6121d233e16b10f7107923c1144a4be9aac7058c12c89a4904369a201

ES0130000001387025DK 7a43f707496a891a3b488efc6a5c19656ed4a3c2bb72ee55d874e43ed8cd2d02
ES0130000001430026QN b273db99783f29e3d95186e3fcce8f5c2f36d0b7536c493d254617cf45375840

*/
