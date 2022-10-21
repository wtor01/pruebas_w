package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"fmt"
	"log"
	"time"
)

/*
*
Procedimiento para la generación de valores iniciales de MVH (billing_measure) en todos los puntos de medida (service_point)
*/
func main() {

	logger := log.Default()

	cnf, _ := config.LoadConfig(logger)

	date := time.Date(2022, 07, 1, 0, 0, 0, 0, cnf.LocalLocation).UTC()
	distributorCode := "0130"
	distributorID := "5140cda8-1daa-4f52-b85d-bf29d45ca62a"
	serviceType := measures.DcServiceType
	pointTypes := []string{"1", "2", "3", "4", "5"}
	meterType := []string{"TLG"}
	for _, pointType := range pointTypes {
		generateInitialMVHByType(date, distributorCode, distributorID, serviceType, pointType, meterType)
	}
	return
	cups := "ES0130000001381064AV"
	generateInitialMVHByCUPS(distributorID, distributorCode, cups, date)
}

/*
*
Método que crea una entrada en billing_measure para que se pueda calcular el MVH de los siguientes meses
*/
func generateInitialMVHByCUPS(distributorID string, distributorCode string, cups string, date time.Time) {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	fmt.Println("Conectando a mongo")
	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	fmt.Println("Conectado a mongo")

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Nuevo inventario")

	inventoryMeasuresRepository := mongo.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	processMeasuresRepository := mongo.NewProcessMeasureRepository(mongoClient, cnf.MeasureDB, cnf.CollectionLoadCurve, cnf.CollectionDailyLoadCurve, cnf.CollectionMonthlyClosure, cnf.CollectionDailyClosure, cnf.LocalLocation)
	billingMeasuresRepository := mongo.NewBillingMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, "billing_measure")

	fmt.Println("Consultando en mongo")

	servicePointConfig, err := inventoryMeasuresRepository.GetMeterConfigByCups(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        cups,
		Time:        date.AddDate(0, 0, -1),
		Distributor: distributorID,
	})

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	fmt.Println("respuesta de mongo ", servicePointConfig)

	magnitudes := servicePointConfig.GetMagnitudesActive()
	billingMeasure := billing_measures.NewBillingMeasure(
		servicePointConfig.Cups(),
		time.Date(2022, 06, 1, 0, 0, 0, 0, cnf.LocalLocation).UTC(),
		date,
		distributorCode,
		distributorID,
		[]measures.PeriodKey{measures.P1, measures.P2, measures.P3},
		magnitudes,
		servicePointConfig.Type,
	)

	billingMeasure.SetContractInfo(servicePointConfig)
	billingMeasure.GenerationDate = time.Now()
	billingMeasure.Status = billing_measures.Billed
	billingMeasure.PointType = servicePointConfig.PointType()
	billingMeasure.RegisterType = servicePointConfig.CurveType

	readingDate, err := getDailyReadingClosureByCups(ctx, processMeasuresRepository, billing_measures.NewProcessMvhEvent(date.AddDate(0, 0, -1), servicePointConfig).Payload)

	if err == nil {

		billingMeasure.SetActualReadingClosure(readingDate)
	}

	billingMeasuresRepository.Save(context.TODO(), billingMeasure)
}

func generateInitialMVHByType(date time.Time, distributorCode string, distributorID string, serviceType measures.ServiceType, pointType string, meterType []string) {

	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	fmt.Println("Conectando a mongo")
	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	fmt.Println("Conectado a mongo")

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Nuevo inventario")

	inventoryMeasuresRepository := mongo.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	processMeasuresRepository := mongo.NewProcessMeasureRepository(mongoClient, cnf.MeasureDB, cnf.CollectionLoadCurve, cnf.CollectionDailyLoadCurve, cnf.CollectionMonthlyClosure, cnf.CollectionDailyClosure, cnf.LocalLocation)
	billingMeasuresRepository := mongo.NewBillingMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, "billing_measure")

	fmt.Println("Consultando en mongo")

	servicePointConfigs, err := inventoryMeasuresRepository.ListMeterConfigByDate(ctx, measures.ListMeterConfigByDateQuery{
		DistributorID: distributorID,
		ServiceType:   serviceType,
		PointType:     pointType,
		MeterType:     meterType,
		Date:          date.AddDate(0, 0, -1),
	})

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	fmt.Println("respuesta de mongo ", len(servicePointConfigs))

	i := 0
	for _, v := range servicePointConfigs {
		if i%100 == 0 {
			fmt.Println("Se han almacenado MVHs: ", i)
		}
		magnitudes := v.GetMagnitudesActive()
		billingMeasure := billing_measures.NewBillingMeasure(
			v.Cups(),
			time.Date(2022, 06, 1, 0, 0, 0, 0, cnf.LocalLocation).UTC(),
			date,
			distributorCode,
			distributorID,
			[]measures.PeriodKey{measures.P1, measures.P2, measures.P3},
			magnitudes,
			v.Type,
		)

		billingMeasure.SetContractInfo(v)
		billingMeasure.GenerationDate = time.Now()
		billingMeasure.Status = billing_measures.Billed
		billingMeasure.PointType = v.PointType()
		billingMeasure.RegisterType = v.CurveType

		readingDate, err := getDailyReadingClosureByCups(ctx, processMeasuresRepository, billing_measures.NewProcessMvhEvent(date.AddDate(0, 0, -1), v).Payload)

		if err == nil {

			billingMeasure.SetActualReadingClosure(readingDate)
		}

		billingMeasuresRepository.Save(context.TODO(), billingMeasure)
		i = i + 1

	}
}

func getDailyReadingClosureByCups(ctx context.Context, processMeasuresRepository process_measures.ProcessedMeasureRepository, dto measures.ProcessMeasurePayload) (measures.DailyReadingClosure, error) {

	monthly, err := processMeasuresRepository.GetMonthlyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date,
	})

	if err == nil {
		return monthly.ToDailyReadingClosure(), nil
	}

	daily, err := processMeasuresRepository.GetProcessedDailyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: dto.MeterConfig.Cups(),
		Date: dto.Date,
	})

	if err != nil {
		return measures.DailyReadingClosure{}, err
	}

	return daily.ToDailyReadingClosure(), nil
}
