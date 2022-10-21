package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/platform/smarkia_api"
	"context"
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

	//process := smarkia.ProcessCurve
	process := smarkia.ProcessClose
	date := time.Date(2022, 6, 1, 0, 0, 0, 0, cnf.LocalLocation).UTC()

	verificarEtapas(process, date)

}

func verificarEtapas(process smarkia.ProcessName, date time.Time) error {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	/*publisherCreator := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	publisher, err := publisherCreator(ctx)*/

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
				if equipment.ID != "2236" {
					continue
				}

				if process == smarkia.ProcessCurve {
					ms, err := smarkiaApi.GetMagnitudes(ctx, equipment.ID, date)
					if err != nil {
						return err
					}

					meter, _ := inventoryRepository.GetMeasureEquipmentBySmarkiaID(ctx, equipment.ID)

					for i := range ms {
						fmt.Println("Medida recibida: ", meter.SerialNumber, "Fecha: ", ms[i].EndDate, " AI: ", ms[i].AI, " AE: ", ms[i].AE)
						ms[i].DistributorCDOS = distributor.CDOS
						ms[i].DistributorID = distributor.ID.String()
						ms[i].MeterSerialNumber = meter.SerialNumber
					}

					return nil

				} else {

					ms, err := smarkiaApi.GetClosinger(ctx, equipment.ID, date.UTC())

					if err != nil {
						return err
					}

					meter, err := inventoryRepository.GetMeasureEquipmentBySmarkiaID(ctx, equipment.ID)

					if err != nil {
						fmt.Println(err)
					}

					for i := range ms {
						ms[i].DistributorCDOS = distributor.CDOS
						ms[i].DistributorID = distributor.ID.String()
						ms[i].MeterSerialNumber = meter.SerialNumber
					}
					return nil

				}
			}
		}
	}
	return nil
}
