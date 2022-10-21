package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cnf config.Config) *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", cnf.DbHost, cnf.DbUser, cnf.DbPass, cnf.DbName, cnf.DbPort)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	db.Logger.LogMode(logger.Info)
	if err != nil {
		panic(err)
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	fmt.Println("Database connection successful.")

	return db
}
