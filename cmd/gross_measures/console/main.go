package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
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

	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	pub, err := publisher(context.Background())

	m := []gross_measures.MeasureCurveWrite{
		{EndDate: time.Date(2022, 02, 07, 00, 45, 00, 000, time.UTC),
			ReadingDate:       time.Date(2022, 05, 07, 00, 45, 00, 000, time.UTC),
			Type:              measures.Incremental,
			DistributorCDOS:   "0130",
			ReadingType:       measures.Curve,
			Status:            measures.Valid,
			Contract:          "",
			MeterSerialNumber: "SAG0155395442",
			ConcentratorID:    "CIR4621936089",
			File:              "",
			DistributorID:     "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
			Origin:            measures.STG,
			Qualifier:         "",
			AI:                0,
			AE:                0,
			R1:                0,
			R2:                0,
			R3:                0,
			R4:                0,
		},
	}
	ev := gross_measures.NewInsertMeasureCurveEvent(m)

	msg, err := json.Marshal(ev)
	attributes := make(map[string]string)
	attributes[event.EventTypeKey] = ev.Type

	if err != nil {
		panic(err)
	}

	err = pub.Publish(context.Background(), cnf.TopicMeasures, msg, attributes)
	if err != nil {
		panic(err)
	}

}
