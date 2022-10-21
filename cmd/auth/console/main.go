package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/firebase"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	firebaseClient, err := firebase.NewOAuthClient()
	if err != nil {
		panic(err)
	}
	db := postgres.New(cnf)
	authRepository := postgres.NewAuthRepository(db)

	svcs := services.NewService(authRepository, firebaseClient)

	err = svcs.CreateUserService.Handler(context.Background(), services.CreateUserDto{
		ID:          uuid.NewV4().String(),
		Email:       "@altostratus.es",
		Name:        "",
		IsAdmin:     true,
		Password:    "123456",
		Permissions: nil,
	})
	fmt.Println(err)
}
