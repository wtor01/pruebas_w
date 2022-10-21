package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/platform/smarkia_api"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

func main() {

	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	//date := time.Date(2022, 7, 3, 0, 0, 0, 0, cnf.LocalLocation)
	distributorID := "5140cda8-1daa-4f52-b85d-bf29d45ca62a"
	cups := "ES0130000000300013JN"
	equipment := "2236"
	ct := "12649"
	//processName := smarkia.ProcessCurve
	processName := smarkia.ProcessClose
	//getCurve(date, distributorID, cups, equipment, ct)
	startDate := time.Date(2022, 6, 1, 0, 0, 0, 0, cnf.LocalLocation)
	endDate := time.Date(2022, 6, 2, 0, 0, 0, 0, cnf.LocalLocation)
	for startDate.Before(endDate) {
		getGrossMeasure(startDate, distributorID, cups, equipment, ct, processName)
		time.Sleep(10 * time.Second)
		startDate = startDate.AddDate(0, 0, 1)
	}
	return

	getGrossMeasures(startDate, endDate, processName)
	return

}

func getGrossMeasure(date time.Time, distributorId string, cups string, equipment string, ct string, processName string) {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	publisherCreator := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	publisher, err := publisherCreator(ctx)

	mongoClient, err := mongo.New(ctx, cnf)

	fmt.Println("Conectado a mongo")

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

	fmt.Println("Conectado a redis")

	db := postgres.New(cnf)

	fmt.Println("Conectado a Postgre")

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(redisClient))

	//smarkiaApi := smarkia_api.NewApi(cnf.SmarkiaToken, cnf.SmarkiaHost, cnf.LocalLocation)

	distributor, err := inventoryRepository.GetDistributorByID(ctx, distributorId)

	//ms, err := smarkiaApi.GetMagnitudes(ctx, equipment, date.UTC())

	generateMsgByCUPSAndSend(ctx, cnf.SmarkiaTopic, distributor, ct, cups, equipment, publisher, date, processName)

}

/*
*
Funcion que inserta en la cola PubSub el mensaje necesario para recuperar de smarkia las curvas de un periodo para todos los CUPS
*/
func getGrossMeasures(initDate time.Time, endDate time.Time, processName string) error {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	publisherCreator := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	publisher, err := publisherCreator(ctx)

	mongoClient, err := mongo.New(ctx, cnf)

	fmt.Println("Conectado a mongo")

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

	fmt.Println("Conectado a redis")

	db := postgres.New(cnf)

	fmt.Println("Conectado a Postgre")

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(redisClient))

	smarkiaApi := smarkia_api.NewApi(cnf.SmarkiaToken, cnf.SmarkiaHost, cnf.LocalLocation)

	distributors, _, err := inventoryRepository.ListDistributors(ctx, inventory.Search{
		Limit: 1000,
		Values: map[string]string{
			"is_smarkia_active": "1",
		},
	})

	for _, distributor := range distributors {
		fmt.Println("Distribuidor: ", distributor.Name)
		cts, err := smarkiaApi.GetCTs(ctx, distributor.SmarkiaId)
		if err != nil {
			return err
		}
		fmt.Println("Recibidos: ", len(cts), " CTs")
		for _, ct := range cts {
			//fmt.Println("CT: ", ct.ID)
			equipments, err := smarkiaApi.GetEquipments(ctx, smarkia.GetEquipmentsQuery{
				Id: ct.ID,
			})
			if err != nil {
				return err
			}
			fmt.Println("Para el CT: ", ct.ID, " Hemos recibido ", len(equipments))
			for _, equipment := range equipments {
				if equipment.CUPS != "ES0130000000300013JN" {
					continue
				}
				startDateAux := initDate
				for startDateAux.Before(endDate) {
					generateMsgAndSend(ctx, cnf.SmarkiaTopic, equipment, ct, distributor, publisher, startDateAux, processName)
					startDateAux = startDateAux.AddDate(0, 0, 1)
				}
				fmt.Println("Equipo: ", equipment.CUPS, " enviado ", endDate)
			}
		}
	}
	return nil
}

func generateMsgAndSend(ctx context.Context, smarkiaTopic string, equipment smarkia.Equipment, ct smarkia.Ct, distributor inventory.Distributor, publisher event.Publisher, date time.Time, processName string) error {

	msg := smarkia.NewEquipmentProcessEvent(smarkia.MessageEquipmentDto{
		ProcessName:     processName,
		DistributorId:   distributor.ID.String(),
		SmarkiaId:       distributor.SmarkiaId,
		CtId:            ct.ID,
		DistributorCDOS: distributor.CDOS,
		Date:            date,
	}, equipment.ID, equipment.CUPS)

	data, err := json.Marshal(msg)
	err = publisher.Publish(ctx, smarkiaTopic, data, msg.GetAttributes())

	return err
}

func generateMsgByCUPSAndSend(ctx context.Context,
	smarkiaTopic string,
	distributor inventory.Distributor,
	ct string, cups string, idEquipment string,
	publisher event.Publisher, date time.Time, processName string) error {

	msg := smarkia.NewEquipmentProcessEvent(smarkia.MessageEquipmentDto{
		ProcessName:     processName,
		DistributorId:   distributor.ID.String(),
		SmarkiaId:       distributor.SmarkiaId,
		CtId:            ct,
		DistributorCDOS: distributor.CDOS,
		Date:            date.AddDate(0, 0, 1),
	}, idEquipment, cups)

	data, err := json.Marshal(msg)
	err = publisher.Publish(ctx, smarkiaTopic, data, msg.GetAttributes())

	fmt.Println("Enviado mensaje para procesar el día", date, " del cups: ", cups)

	return err
}
