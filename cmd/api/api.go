package main

import (
	"bitbucket.org/sercide/data-ingestion/cmd/aggregations/http/aggregations"
	aggregation_config "bitbucket.org/sercide/data-ingestion/cmd/aggregations/http/config"
	"bitbucket.org/sercide/data-ingestion/cmd/aggregations/http/features"
	auth_http "bitbucket.org/sercide/data-ingestion/cmd/auth/http"
	billing_measures_dashboard "bitbucket.org/sercide/data-ingestion/cmd/billing_measures/http/dashboard"
	billing_measures_scheduler "bitbucket.org/sercide/data-ingestion/cmd/billing_measures/http/scheduler"
	"bitbucket.org/sercide/data-ingestion/cmd/billing_measures/http/self_consumption"
	gross_measures_dashboard "bitbucket.org/sercide/data-ingestion/cmd/gross_measures/http/dashboard"
	stats "bitbucket.org/sercide/data-ingestion/cmd/gross_measures/http/dashboard_stats"
	smarkia_http "bitbucket.org/sercide/data-ingestion/cmd/gross_measures/http/smarkia"
	"bitbucket.org/sercide/data-ingestion/cmd/gross_measures/http/supply_point"
	inventory_http "bitbucket.org/sercide/data-ingestion/cmd/inventory/http"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/calendar"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/festive_days"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/generate_calendars"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/geographic"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/seasons"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/tariff"
	process_measures_closure "bitbucket.org/sercide/data-ingestion/cmd/process_measures/http/closures"
	process_measures_dashboard "bitbucket.org/sercide/data-ingestion/cmd/process_measures/http/dashboard"
	process_measure_dashboard_stats_http "bitbucket.org/sercide/data-ingestion/cmd/process_measures/http/dashboard_stats"
	process_measure_scheduler_http "bitbucket.org/sercide/data-ingestion/cmd/process_measures/http/scheduler"
	validations "bitbucket.org/sercide/data-ingestion/cmd/validations/http/validations"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redisRepos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"log"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, s := server.NewHttpServer(context.Background(), "", cnf.Port, cnf.ShutdownTimeout)

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
	authMiddleware := middleware.NewAuth(cnf)
	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redisRepos.New(cnf)
		defer redisClient.Close()
	}

	tp := telemetry.NewTracerProvider("api")
	defer tp.ForceFlush(ctx)

	auth_http.Register(&s, authMiddleware)
	gross_measures_dashboard.Register(&s, authMiddleware, db, cnf, mongoClient, redisClient)
	validations.Register(&s, authMiddleware, db, mongoClient, redisClient, cnf)
	inventory_http.Register(&s, authMiddleware, db, redisClient, mongoClient, cnf)
	process_measure_scheduler_http.Register(&s, cnf, authMiddleware, db)
	process_measures_dashboard.Register(&s, authMiddleware, db, cnf, mongoClient, redisClient)
	process_measure_dashboard_stats_http.Register(&s, cnf, authMiddleware, mongoClient)
	billing_measures_scheduler.Register(&s, cnf, authMiddleware, db)
	geographic.Register(&s, cnf, authMiddleware, db, redisClient)
	calendar.Register(&s, cnf, authMiddleware, db, redisClient)
	tariff.Register(&s, authMiddleware, db, redisClient)

	billing_measures_dashboard.Register(&s, cnf, authMiddleware, mongoClient)

	process_measures_closure.Register(&s, authMiddleware, mongoClient, cnf, db, redisClient)

	aggregation_config.Register(&s, authMiddleware, db, cnf)
	aggregations.Register(&s, cnf, authMiddleware, mongoClient)
	features.Register(&s, cnf, authMiddleware, db)

	seasons.Register(&s, authMiddleware, db, redisClient)
	festive_days.Register(&s, cnf, authMiddleware, db, redisClient)
	generate_calendars.Register(&s, authMiddleware, db, redisClient)

	self_consumption.Register(&s, authMiddleware, mongoClient, cnf)
	smarkia_http.Register(&s, authMiddleware, db, cnf, mongoClient, redisClient)
	stats.Register(&s, cnf, authMiddleware, mongoClient)
	supply_point.Register(&s, authMiddleware, db, cnf, mongoClient, redisClient)
	
	s.Register(func(router *gin.RouterGroup) {
		server.ServeOpenapi(router, server.SwaggerUIOpts{
			SpecURL: "docs/api.yaml",
		})
	})

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
