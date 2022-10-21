package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"fmt"
	"log"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}
	db := postgres.New(cnf)

	//postgres.Rollback(db)

	err = postgres.Migrate(db)
	if err != nil {
		logger.Fatal(err)
	}
	//postgres.CreateDistributor(db)
	//postgres.CreateEquipments(db)

	fmt.Println(err)
}
