package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	ProcessCurveEventType          = "PROCESS_MEASURES/CURVE"
	ProcessBillingClosureEventType = "PROCESS_MEASURES/BILLING_CLOSURE"
	ProcessDailyClosureEventType   = "PROCESS_MEASURES/DAILY_CLOSURE"
	SchedulerEventType             = "PROCESS_MEASURES/INIT"
	TypeDistributorProcess         = "PROCESS/DISTRIBUTOR"
)

func main() {

	logger := log.Default()

	cnf, err := config.LoadConfig(logger)

	/* VARIABLES */
	distributorId := "5140cda8-1daa-4f52-b85d-bf29d45ca62a"
	//cups := "ES0130000001392019CR"
	startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, cnf.LocalLocation)
	endDate := time.Date(2022, 8, 1, 0, 0, 0, 0, cnf.LocalLocation)
	//ServiceTypes := []measures.ServiceType{measures.DcServiceType, measures.GdServiceType, measures.DdServiceType}
	ServiceTypes := []measures.ServiceType{measures.DcServiceType, measures.GdServiceType}
	PointTypes := []string{"1", "2", "3", "4", "5"}
	//PointTypes := []string{"4"}
	//MeterTypes := []string{"TLG", "TLM", "OTHER"}
	MeterTypes := []string{"TLM"}
	//MeterTypes := []string{"TLG"}
	ReadingTypes := []measures.ReadingType{measures.Curve}
	/* VARIABLES */

	cups := "ES0130000000300013JN"
	for startDate.Before(endDate) {
		processCUPS(distributorId, cups, startDate, measures.Curve)
		startDate = startDate.AddDate(0, 0, 1)
	}
	return

	for _, measureType := range ReadingTypes {
		startDateAux := startDate
		for startDateAux.Before(endDate) {
			for _, serviceType := range ServiceTypes {
				for _, pointType := range PointTypes {
					for _, meterType := range MeterTypes {
						processAllCups(distributorId, startDateAux, measureType, serviceType, pointType, meterType)
					}
				}
			}
			startDateAux = startDateAux.AddDate(0, 0, 1)
		}
	}

	fmt.Println("Fin del proceso")
	fmt.Println(err)

	return

	//simulateProcessDistributorStep(distributorId, startDate, measures.Curve, measures.DcServiceType, "4", "TLG")
	//cups := "ES0130000000344035HN"
	//cups := "ES0130000001381064AV"
	//cups := "ES0130000001348006SX"
	//date := time.Date(2022, 10, 1, 0, 0, 0, 0, cnf.LocalLocation)
	/*simulateProcessMonthlyClosure(distributorId, cups, date)
	return*/
	//[END] Procesar todos los cups para un rango de días
	//[START] Procesar un cups para un rango de días
	//for startDate.Before(endDate) {
	//processCUPS(distributorId, cups, date, measures.DailyClosure)
	//startDate = startDate.AddDate(0, 0, 1)
	//}
	return
	//[START] Procesar todos los cups para un rango de días
	//for _, measureType := range measures.ReadingTypes {

}

func simulateProcessMonthlyClosure(distributorID string, CUPS string, date time.Time) {
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

	processRepository := mongo.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	grossRepository := mongo.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)

	ms, err := grossRepository.ListDailyCloseMeasures(ctx, gross_measures.QueryListForProcessClose{
		ReadingType:  measures.BillingClosure,
		SerialNumber: "83020034",
		Date:         date,
	})

	if err != nil {
		fmt.Println("ERROR en consulta de medida: ", err)
	}

	if len(ms) == 0 {
		fmt.Println("ERROR no hay resultados de medida: ", err)
	}

	processedMeasureLM, err := processRepository.GetMonthlyClosureByCup(ctx, process_measures.QueryClosedCupsMeasureOnDate{
		CUPS: CUPS,
		Date: date.AddDate(0, -1, 0),
	})

	if err != nil {
		fmt.Println("ERROR: ", err)
		//return
	}

	bestIncrementClose := gross_measures.MeasureCloseWrite{}
	bestAbsoluteClose := gross_measures.MeasureCloseWrite{}

	//Escoger las mejores medidas de incremento y absolutas para ese meter Id
	for _, measure := range ms {

		fmt.Println("La medida del cierre es:", measure.Type)
		fmt.Println("Actualmente el cierre absoluto es: ", bestAbsoluteClose.Type)
		fmt.Println("Actualmente el cierre incremental es: ", bestIncrementClose.Type)

		if measure.Type == measures.Absolute && bestAbsoluteClose.Type == "" {
			bestAbsoluteClose = measure
		}

		if measure.Type == measures.Incremental && bestIncrementClose.Type == "" {
			bestIncrementClose = measure
		}

		if bestIncrementClose.Type != "" && bestAbsoluteClose.Type != "" {
			break
		}

	}

	if bestIncrementClose.Type == "" && bestAbsoluteClose.Type == "" {
		fmt.Println("invalid measures")
		return
	}

	calendarPeriods := GetProcessedMonthlyClosureCalendar(bestAbsoluteClose, bestIncrementClose)

	fmt.Println("Termina OK: ", calendarPeriods)
	return
	fmt.Println("ResultadoLastMesure: ", processedMeasureLM)
}

func simulateProcessDistributorStep(distributorId string, date time.Time, processType measures.ReadingType, serviceType measures.ServiceType, pointType string, meterType string) {
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

	inventoryMeasuresRepository := mongo.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	meterConfigs, err := inventoryMeasuresRepository.ListMeterConfigByDate(ctx, measures.ListMeterConfigByDateQuery{
		DistributorID: distributorId,
		ServiceType:   serviceType,
		PointType:     pointType,
		MeterType:     []string{meterType},
		ReadingType:   processType,
		Date:          date,
		Limit:         10,
		Offset:        0,
	})

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Datos recibidos: ", len(meterConfigs))

	for _, meterConfig := range meterConfigs {
		msg := getMessageAttributes(processType, meterConfig, date)
		fmt.Println("Dia ", date, " Mensaje ", msg)
	}
}

/**
 * Función que envía a la cola de procesamiento de medida un mensaje para que procese la medida de un tipo para un cups en un día concreto
 */
func processCUPS(distributorId string, cups string, date time.Time, processType measures.ReadingType) {
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
		Distributor: distributorId,
		CUPS:        cups,
		Time:        date,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := getMessageAttributes(processType, meterConfig, date)

	fmt.Println("Dia ", date, " CUPS ", cups, " tipo ", processType)
	for _, msgIterator := range msg {
		data, _ := json.Marshal(msgIterator)
		err = publisher.Publish(ctx, cnf.ProcessedMeasureTopic, data, msgIterator.GetAttributes())
	}

}

/**
 * Función que envía a la cola de procesamiento de medida un mensaje para que procese toda la medida de un distribuidor
 */
func processAllCups(distributorId string, date time.Time, processType measures.ReadingType, serviceType measures.ServiceType, pointType string, meterType string) {
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
		"type": SchedulerEventType,
	}

	msg := measures.NewSchedulerEvent(process_measures.SchedulerEventType, measures.SchedulerEventPayload{
		ID:            "0bc22922-22bf-11ed-a4bd-42706f29bceb",
		DistributorId: distributorId,
		Name:          "ProcesamientoResumen",
		Description:   "Procesamiento de medida para resumenes diarios",
		ServiceType:   serviceType,
		PointType:     pointType,
		MeterType:     []string{meterType},
		ReadingType:   processType,
		Format:        "* * 1 * *",
		Date:          date,
	})

	fmt.Println("Dia: ", date, " Tipo procesamiento: ", processType, " Servicio: ", serviceType, " Tipo de punto: ", pointType, " Tipo de Contador: ", meterType)

	data, err := json.Marshal(msg)

	if err != nil {
		logger.Fatal(err)
	}

	err = publisher.Publish(ctx, cnf.ProcessedMeasureTopic, data, attributes)

	if err != nil {
		logger.Fatal(err)
	}
}

func getMessageAttributes(processType measures.ReadingType, meterConfig measures.MeterConfig, date time.Time) []measures.ProcessMeasureEvent {
	var msg []measures.ProcessMeasureEvent

	switch processType {
	case measures.Curve:
		{
			if utils.InSlice(meterConfig.CurveType, []measures.RegisterType{measures.Hourly, measures.Both}) {
				msg = append(msg, process_measures.NewProcessCurveEvent(date, meterConfig, measures.HourlyMeasureCurveReadingType))
			}

			if utils.InSlice(meterConfig.CurveType, []measures.RegisterType{measures.QuarterHour, measures.Both}) {
				msg = append(msg, process_measures.NewProcessCurveEvent(date, meterConfig, measures.QuarterMeasureCurveReadingType))
			}

		}
	case measures.BillingClosure:
		{
			msg = append(msg, process_measures.NewProcessBillingClosureEvent(date, meterConfig))
		}
	case measures.DailyClosure:
		{
			msg = append(msg, process_measures.NewProcessDailyClosureEvent(date, meterConfig))
		}
	}

	return msg
}

func GetProcessedMonthlyClosureCalendar(bestAbsoluteMeasure, bestIncrementClose gross_measures.MeasureCloseWrite) process_measures.ProcessedMonthlyClosureCalendar {

	var calendarPeriods process_measures.ProcessedMonthlyClosureCalendar

	for _, p := range append(measures.ValidPeriodsCurve, measures.P0) {
		switch p {
		case measures.P0:
			calendarPeriods.P0 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P1:
			calendarPeriods.P1 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P2:
			calendarPeriods.P2 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P3:
			calendarPeriods.P3 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P4:
			calendarPeriods.P4 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P5:
			calendarPeriods.P5 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		case measures.P6:
			calendarPeriods.P6 = &process_measures.ProcessedMonthlyClosurePeriod{
				ProcessedDailyClosurePeriod: process_measures.ProcessedDailyClosurePeriod{
					Filled:           true,
					ValidationStatus: measures.Invalid,
				},
			}
		}
	}

	for _, period := range bestAbsoluteMeasure.Periods {

		var processedMonthlyClosurePeriod *process_measures.ProcessedMonthlyClosurePeriod

		switch period.Period {
		case measures.P0:
			processedMonthlyClosurePeriod = calendarPeriods.P0
		case measures.P1:
			processedMonthlyClosurePeriod = calendarPeriods.P1
		case measures.P2:
			processedMonthlyClosurePeriod = calendarPeriods.P2
		case measures.P3:
			processedMonthlyClosurePeriod = calendarPeriods.P3
		case measures.P4:
			processedMonthlyClosurePeriod = calendarPeriods.P4
		case measures.P5:
			processedMonthlyClosurePeriod = calendarPeriods.P5
		case measures.P6:
			processedMonthlyClosurePeriod = calendarPeriods.P6
		}

		if processedMonthlyClosurePeriod == nil {
			continue
		}
		processedMonthlyClosurePeriod.Values = measures.Values{
			AI: period.AI,
			AE: period.AE,
			R1: period.R1,
			R2: period.R2,
			R3: period.R3,
			R4: period.R4,
		}
		processedMonthlyClosurePeriod.ValidationStatus = measures.Valid
		processedMonthlyClosurePeriod.Filled = false
	}

	for _, period := range bestIncrementClose.Periods {
		var calendarPeriod *process_measures.ProcessedMonthlyClosurePeriod
		switch period.Period {
		case measures.P0:
			calendarPeriod = calendarPeriods.P0
		case measures.P1:
			calendarPeriod = calendarPeriods.P1
		case measures.P2:
			calendarPeriod = calendarPeriods.P2
		case measures.P3:
			calendarPeriod = calendarPeriods.P3
		case measures.P4:
			calendarPeriod = calendarPeriods.P4
		case measures.P5:
			calendarPeriod = calendarPeriods.P5
		case measures.P6:
			calendarPeriod = calendarPeriods.P6
		}
		if calendarPeriod == nil {
			continue
		}
		calendarPeriod.AIi = period.AI
		calendarPeriod.AEi = period.AE
		calendarPeriod.R1i = period.R1
		calendarPeriod.R2i = period.R2
		calendarPeriod.R3i = period.R3
		calendarPeriod.R4i = period.R4
	}

	return calendarPeriods
}
