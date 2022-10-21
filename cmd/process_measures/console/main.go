package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"encoding/json"
	"log"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("Europe/Madrid")
	date := time.Date(2022, 8, 31, 0, 0, 0, 0, loc)
	endDate := time.Date(2022, 8, 32, 0, 0, 0, 0, loc)
	for date.Before(endDate) {
		processMeasureByMeter("ES0130000000361013CQ", "5140cda8-1daa-4f52-b85d-bf29d45ca62a", date)
		date = date.AddDate(0, 0, 1)
	}
	/*logger := log.Default()
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
	repositoryFestiveDays := redis_repos.NewProcessMeasureFestiveDays(redisClient)
	processRepository := mongo_repos.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)
	grossRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	grossMeasuresClient := gross_measures_client.NewClient(services2.NewListDailyCurveMeasures(grossRepository), services2.NewListDailyCloseMeasures(grossRepository))
	svc := services.NewProcessCurve(
		grossMeasuresClient,
		processRepository,
		nil,
		repositoryFestiveDays,
		cnf.LocalLocation,
	)
	err = svc.Handle(context.Background(), measures.ProcessMeasurePayload{
		MeterId:           "80c50b06-28cd-49d6-97e8-079eed2c5de7",
		Date:              time.Date(2022, 6, 1, 0, 0, 0, 0, time.Local).UTC(),
		DistributorID:     "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
		DistributorCode:   "0130",
		CUPS:              "ES0130000000357054DJ",
		MeterSerialNumber: "",
		MeasurePointType:  "",
		PriorityContract:    "1",
		MeterConfigId:     "",
		ServiceType:       "D-C",
		PointType:         "5",
		RecoverMeasuresId: "",
		CurveType:         string(measures.Hourly),
		TariffID:          "2.0TD",
		CalendarCode:      "TD3",
		Coefficient:       "",
		P1Demand:          0,
		P2Demand:          0,
		P3Demand:          0,
		P4Demand:          0,
		P5Demand:          0,
		P6Demand:          0,
		GeographicID:      "f719d66f-d3c5-420e-9120-6a721e86a8e0",
	})
	fmt.Println(err)*/
}

func list() {
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

	measureRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	measureRepository.ListDailyCurveMeasures(context.Background(), gross_measures.QueryListForProcessCurve{
		SerialNumber: "85370c01-fcec-4513-b651-33d62a1117ed",
		Date:         time.Date(2022, 04, 30, 0, 0, 0, 0, time.Local),
	})
}

func processMeasures() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)

	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	publisherCreator := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	publisher, err := publisherCreator(ctx)
	if err != nil {
		panic(err)
	}

	loc, _ := time.LoadLocation(cnf.TimeZone)
	startDate := time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC()
	endDate := time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC()

	attributes := map[string]string{
		"type": "PROCESS_MEASURES/INIT",
	}

	types := map[string]string{
		"curve": "",
		"close": "",
	}

	for startDate.Before(endDate) {
		for t := range types {
			msg := measures.NewSchedulerEvent(process_measures.SchedulerEventType, measures.SchedulerEventPayload{
				ID:            "7a4ab092-e2b1-11ec-a897-b63aebc103a9",
				DistributorId: "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Name:          t,
				Description:   "ABCD",
				ServiceType:   "D-C",
				PointType:     "5",
				MeterType:     []string{"TLM"},
				ReadingType:   measures.ReadingType(t),
				Format:        "* * 1 * *",
				Date:          startDate,
			})

			data, err := json.Marshal(msg)

			if err != nil {
				logger.Fatal(err)
			}

			err = publisher.Publish(ctx, cnf.SmarkiaTopic, data, attributes)

			if err != nil {
				logger.Fatal(err)
			}
		}

		startDate = startDate.AddDate(0, 0, 1)
	}
}

func processMeasureByMeter(cups, distributor string, date time.Time) {
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

	publisherCreator := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	publisher, err := publisherCreator(ctx)
	if err != nil {
		panic(err)
	}

	inventoryRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	meter, err := inventoryRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        cups,
		Time:        date,
		Distributor: distributor,
	})

	if err != nil {
		return
	}

	event := measures.NewProcessMeasureEvent(
		process_measures.ProcessCurveEventType,
		date,
		meter,
	)
	data, err := json.Marshal(event)

	attributes := map[string]string{
		"type": process_measures.ProcessCurveEventType,
	}

	if err != nil {
		logger.Fatal(err)
	}

	err = publisher.Publish(ctx, cnf.ProcessedMeasureTopic, data, attributes)

	logger.Println("Dia: ", date.Format("2006/01/02"), " - Cups: ", cups, " - Distributor: ", distributor)
	if err != nil {
		logger.Fatal(err)
	}

}
