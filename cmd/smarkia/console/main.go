package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	"context"
	"encoding/json"
	"log"
	"time"
)

func main() {
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
	startDate := time.Date(2022, 7, 30, 0, 0, 0, 0, loc).UTC()
	endDate := time.Date(2022, 8, 2, 0, 0, 0, 0, loc).UTC()

	attributes := map[string]string{
		"type": "SMARKIA/REQUEST",
	}

	type MessagePayload struct {
		Type    string `json:"type"`
		Payload struct {
			ProcessName string    `json:"process_name"`
			Date        time.Time `json:"date"`
		} `json:"payload"`
	}

	types := map[string]string{
		"curve": "SMARKIA/FETCH_CURVE",
		"close": "SMARKIA/FETCH_CLOSE",
	}

	for startDate.Before(endDate) {

		for t, v := range types {
			msg := MessagePayload{
				Type: v,
				Payload: struct {
					ProcessName string    `json:"process_name"`
					Date        time.Time `json:"date"`
				}{
					ProcessName: t,
					Date:        startDate,
				},
			}

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
