package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"sort"
	"time"

	mongodb "go.mongodb.org/mongo-driver/mongo"
)

func main() {

	logger := log.Default()

	cnf, _ := config.LoadConfig(logger)

	//cups := "ES0130000001381064AV"
	//cups := "ES0130000001395045JZ"
	//cups := "ES0130000000312025YH"
	//cups := "ES0130000000300026ZW"
	//cups := "ES0130000000361058ES"
	//cups := "ES0130000000319044NE"
	cups := "ES0130000000300013JN"

	date := time.Date(2022, 7, 31, 0, 0, 0, 0, cnf.LocalLocation).UTC()
	distributorId := "5140cda8-1daa-4f52-b85d-bf29d45ca62a"
	ServiceTypes := []measures.ServiceType{measures.DcServiceType, measures.GdServiceType}
	PointTypes := []string{"1", "2", "3", "4", "5"}
	MeterTypes := []string{"TLM"}
	//serviceType := measures.DcServiceType
	//pointType := "5"
	//meterType := "TLG"

	executeMVH(distributorId, cups, date)
	//simulateMVH(distributorId, cups, date)
	return

	for _, serviceType := range ServiceTypes {
		for _, pointType := range PointTypes {
			for _, meterType := range MeterTypes {
				executeAll(distributorId, date, serviceType, pointType, meterType)
			}
		}
	}
	return

}

func executeAll(distributorId string, date time.Time, serviceType measures.ServiceType, pointType string, meterType string) {

	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	fmt.Println(cnf.LocalLocation)

	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	publisherCreator := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	publisher, err := publisherCreator(ctx)
	if err != nil {
		panic(err)
	}

	attributes := map[string]string{
		"type": billing_measures.SchedulerEventType,
	}

	msg := measures.NewSchedulerEvent(billing_measures.SchedulerEventType, measures.SchedulerEventPayload{
		ID:            "0bc22922-22bf-11ed-a4bd-42706f29bceb",
		DistributorId: distributorId,
		Name:          "EjecucionMVH",
		Description:   "",
		ServiceType:   serviceType,
		PointType:     pointType,
		MeterType:     []string{meterType},
		Format:        "* * 1 * *",
		Date:          date,
	})

	fmt.Println("Dia: ", date, " Servicio: ", serviceType, " Tipo de punto: ", pointType, " Tipo de Contador: ", meterType)

	data, err := json.Marshal(msg)

	if err != nil {
		logger.Fatal(err)
	}

	err = publisher.Publish(ctx, cnf.TopicBillingMeasures, data, attributes)

	if err != nil {
		logger.Fatal(err)
	}
}

/*
*
Función que envía a la cola de procesamiento del MVH un cups y una fecha para hacer el cálculo del MVH
*/
func executeMVH(distributorId string, cups string, date time.Time) {
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

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	inventoryMeasuresRepository := mongo.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	meterConfig, err := inventoryMeasuresRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        cups,
		Time:        date,
		Distributor: distributorId,
	})

	msg := billing_measures.NewProcessMvhEvent(date, meterConfig)

	data, err := json.Marshal(msg)

	if err != nil {
		fmt.Println("Error en CUP:", meterConfig.Cups())
		return
	}
	attributes := map[string]string{"type": "BILLING_MEASURES/PROCESS_MVH"}
	err = publisher.Publish(ctx, cnf.TopicBillingMeasures, data, attributes)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

/*
*
Función que simula la ejecución del MVH
*/
func simulateMVH(distributorId string, cups string, date time.Time) {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := postgres.New(cnf)
	var redisClient *redis.Client
	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}

	inventoryMeasuresRepository := mongo.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	billingMeasuresRepository := mongo.NewBillingMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, "billing_measure")
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(redisClient))
	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))
	repositoryFestiveDays := redis_repos.NewProcessMeasureFestiveDays(redisClient)

	meterConfig, err := inventoryMeasuresRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        cups,
		Time:        date,
		Distributor: distributorId,
	})

	dto := measures.ProcessMeasurePayload{
		MeterConfig: meterConfig,
		Date:        date,
	}

	tariff, err := masterTablesClient.GetTariff(ctx, clients.GetTariffDto{
		ID: meterConfig.TariffID(),
	})

	fmt.Println("tarifa: ", tariff, "calendar: ", tariff.CalendarId)

	calendarPeriodDay, err := repositoryFestiveDays.GetCalendarPeriod(ctx, measures.SearchCalendarPeriod{
		Day:          dto.Date,
		GeographicID: tariff.GeographicId,
		CalendarCode: tariff.CalendarId,
		Location:     cnf.LocalLocation,
	})

	periods := calendarPeriodDay.GetAllPeriods()

	sort.Slice(periods, func(i, j int) bool {
		return periods[i] > periods[j]
	})

	magnitudes := meterConfig.GetMagnitudesActive()

	lastBillingMeasure, err := billingMeasuresRepository.Last(ctx, billing_measures.QueryLast{
		CUPS: meterConfig.Cups(),
	})

	billingMeasure := billing_measures.NewBillingMeasure(
		cups,
		lastBillingMeasure.EndDate,
		date.AddDate(0, 0, 1),
		lastBillingMeasure.DistributorCode,
		lastBillingMeasure.DistributorID,
		periods,
		magnitudes,
		meterConfig.Type,
	)

	if err != nil || lastBillingMeasure.EndDate.IsZero() {
		fmt.Println("Supervision: No hay billing measure anterior")
	}

	previous, _ := billingMeasuresRepository.GetPrevious(ctx, billing_measures.GetPrevious{
		CUPS:     billingMeasure.CUPS,
		InitDate: billingMeasure.InitDate,
		EndDate:  billingMeasure.EndDate,
	})

	billingMeasure.CalculateVersionByPreviousBillingMeasure(previous)

	billingMeasure.SetContractInfo(meterConfig)
	billingMeasure.SetCoefficient(tariff)

	billingMeasure.Status = billing_measures.Calculating

	billingMeasure.Inaccessible = meterConfig.Inaccessible
	billingMeasure.PointType = meterConfig.PointType()
	billingMeasure.RegisterType = meterConfig.CurveType

	err = setBillingLoadCurve(ctx, dto, &lastBillingMeasure, &billingMeasure, mongoClient, cnf)

	if err != nil {
		fmt.Println("Debería ser supervisión, pero no es y aqui peta con TLM")
	}

	err = setReadingsClosure(ctx, dto, &lastBillingMeasure, &billingMeasure, mongoClient, cnf)

	if err != nil {
		fmt.Println("Peta aqui")
	}

	billingMeasure.CalcAtrBalance()

	billingMeasure.CalcAtrVsCurve()

}

func getDailyReadingClosureByCups(ctx context.Context, dto measures.ProcessMeasurePayload, mongoClient *mongodb.Client,
	cnf config.Config) (measures.DailyReadingClosure, error) {

	processedMeasureRepository := mongo.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)
	monthly, err := processedMeasureRepository.GetMonthlyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date,
	})

	if err == nil {
		return monthly.ToDailyReadingClosure(), nil
	}

	daily, err := processedMeasureRepository.GetProcessedDailyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date,
	})

	if err != nil {
		return measures.DailyReadingClosure{}, err
	}

	return daily.ToDailyReadingClosure(), nil
}

func setReadingsClosure(
	ctx context.Context,
	dto measures.ProcessMeasurePayload,
	lastBillingMeasure *billing_measures.BillingMeasure,
	billingMeasure *billing_measures.BillingMeasure,
	mongoClient *mongodb.Client,
	cnf config.Config,
) error {

	billingMeasure.SetPreviousReadingClosure(lastBillingMeasure.ActualReadingClosure)

	readingDate, err := getDailyReadingClosureByCups(ctx, dto, mongoClient, cnf)

	if err == nil {
		if !readingDate.InitDate.IsZero() && readingDate.InitDate.Before(lastBillingMeasure.EndDate) {
			fmt.Println("Supervision: error fechas cierre")
		}
		billingMeasure.SetActualReadingClosure(readingDate)
	}

	return nil
}

func setBillingLoadCurve(
	ctx context.Context,
	dto measures.ProcessMeasurePayload,
	lastBillingMeasure *billing_measures.BillingMeasure,
	b *billing_measures.BillingMeasure,
	mongoClient *mongodb.Client,
	cnf config.Config,
) error {

	processedMeasureRepository := mongo.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	processedLoadCurve, err := processedMeasureRepository.ProcessedLoadCurveByCups(ctx, process_measures.QueryProcessedLoadCurveByCups{
		CUPS:      dto.MeterConfig.Cups(),
		StartDate: lastBillingMeasure.EndDate,
		EndDate:   dto.Date.AddDate(0, 0, 1),
		Status:    measures.Valid,
	})

	if err != nil {
		return err
	}

	b.SetBillingLoadCurve(utils.MapSlice(processedLoadCurve, func(item process_measures.ProcessedLoadCurve) billing_measures.BillingLoadCurve {
		return billing_measures.BillingLoadCurve{
			EndDate:          item.EndDate,
			Origin:           item.Origin,
			AI:               item.AI,
			AE:               item.AE,
			R1:               item.R1,
			R2:               item.R2,
			R3:               item.R3,
			R4:               item.R4,
			Period:           item.Period,
			MeasurePointType: item.MeasurePointType,
		}
	}))

	return nil
}
