package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	"context"
	"encoding/csv"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type InsertConsumProfileDTO struct {
	File string
}

type InsertConsumProfile struct {
	storageCreator storage.StorageCreator
	repository     billing_measures.ConsumProfileRepository
}

func NewInsertConsumProfile(storageCreator storage.StorageCreator, repository billing_measures.ConsumProfileRepository) *InsertConsumProfile {
	return &InsertConsumProfile{storageCreator: storageCreator, repository: repository}
}

func (svc InsertConsumProfile) Handler(ctx context.Context, dto InsertConsumProfileDTO) error {

	sto, err := svc.storageCreator(ctx)

	if err != nil {
		return err
	}

	reader, err := sto.NewReader(ctx, dto.File)
	if err != nil {
		return err
	}
	csvReader := csv.NewReader(reader)
	var errorsSave uint64 = 0
	var gorutinesToSave int64 = 0

	var wg sync.WaitGroup
	for {
		if gorutinesToSave > 100 {
			wg.Wait()
		}
		row, err := csvReader.Read()
		if err != nil {
			break
		}
		consumProfile := billing_measures.ConsumProfile{}
		if len(row) != 8 {
			return errors.New("error invalid header file")
		}

		date, err := time.Parse("2006-01-02 15:04:05.000", row[0])

		if err != nil {
			continue
		}

		gmt, err := strconv.Atoi(row[1])

		if err != nil {
			continue
		}

		consumProfile.Date = date.In(time.UTC).Add(-(time.Hour * time.Duration(gmt)))

		version, err := strconv.Atoi(row[2])

		if err != nil {
			continue
		}
		consumProfile.Version = version

		consumProfile.Type = billing_measures.ConsumProfileType(row[3])

		if coefA, err := strconv.ParseFloat(row[4], 64); err == nil {
			consumProfile.CoefA = &coefA
		}
		if coefB, err := strconv.ParseFloat(row[5], 64); err == nil {
			consumProfile.CoefB = &coefB
		}
		if coefC, err := strconv.ParseFloat(row[6], 64); err == nil {
			consumProfile.CoefC = &coefC
		}
		if coefD, err := strconv.ParseFloat(row[7], 64); err == nil {
			consumProfile.CoefD = &coefD
		}

		wg.Add(1)
		go func(cc *billing_measures.ConsumProfile) {
			atomic.AddInt64(&gorutinesToSave, 1)
			defer wg.Done()
			defer func() {
				atomic.AddInt64(&gorutinesToSave, -1)
			}()
			err = svc.repository.Save(ctx, *cc)
			if err != nil {
				atomic.AddUint64(&errorsSave, 1)
			}
		}(&consumProfile)
	}
	wg.Wait()

	if errorsSave >= 1 {
		return errors.New("error saving ConsumProfile from file")
	}

	return nil
}
