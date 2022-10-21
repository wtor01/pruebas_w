package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/parsers"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type ParseFilesService struct {
	publisher      event.PublisherCreator
	topicMeasure   string
	handlers       []parsers.ParserCurveI
	handlersClose  []parsers.ParserCloseI
	storage        storage.StorageCreator
	storageFail    string
	storageSuccess string
}

func NewParseFilesService(publisher event.PublisherCreator, topicMeasure string, storage storage.StorageCreator, storageFail string, storageSuccess string, loc *time.Location) ParseFilesService {

	parserOptions := parsers.NewParser(storage, loc)

	return ParseFilesService{
		publisher:    publisher,
		topicMeasure: topicMeasure,
		handlers: []parsers.ParserCurveI{
			parsers.NewPrimeS02(parserOptions),
			parsers.NewCurveTpl(parserOptions),
			parsers.NewCurvaCsv(parserOptions),
		},
		handlersClose: []parsers.ParserCloseI{
			parsers.NewCloseCsv(parserOptions),
			parsers.NewResumenCsv(parserOptions),
			parsers.NewClosingTpl(parserOptions),
			parsers.NewPrimeS04(parserOptions),
			parsers.NewPrimeS05(parserOptions),
			parsers.NewReadingsTpl(parserOptions),
		},
		storage:        storage,
		storageFail:    storageFail,
		storageSuccess: storageSuccess,
	}
}

func (svc ParseFilesService) Handle(ctx context.Context, dto gross_measures.HandleFileDTO) error {
	handler := svc.getHandler(ctx, dto)
	if handler != nil {
		return svc.HandleFile(ctx, dto)
	}
	return svc.HandleFileClose(ctx, dto)
}

func (svc ParseFilesService) HandleFile(ctx context.Context, dto gross_measures.HandleFileDTO) error {
	if !svc.isValidFile(dto) {
		return nil
	}

	handler := svc.getHandler(ctx, dto)
	if handler == nil {
		return errors.New("not handler")
	}

	var err error

	defer svc.moveFileToDone(ctx, dto, err)

	measureList, err := handler.Handle(ctx, dto)
	if err != nil {
		return err
	}

	events := gross_measures.ListMeasureCurveWriteToEvents(measureList, gross_measures.MaxMeasuresInEvent)

	err = event.PublishAllEvents(ctx, svc.topicMeasure, svc.publisher, events)

	return err
}

func (svc ParseFilesService) HandleFileClose(ctx context.Context, dto gross_measures.HandleFileDTO) error {
	if !svc.isValidFile(dto) {
		return nil
	}

	handler := svc.getHandlerClose(ctx, dto)
	if handler == nil {
		return errors.New("not handler")
	}
	var err error

	defer svc.moveFileToDone(ctx, dto, err)

	measureList, err := handler.Handle(ctx, dto)

	if err != nil {
		return err
	}

	events := gross_measures.ListMeasureCloseWriteToEvents(measureList, gross_measures.MaxMeasuresInEvent)

	err = event.PublishAllEvents(ctx, svc.topicMeasure, svc.publisher, events)

	return err
}

func (svc ParseFilesService) isValidFile(dto gross_measures.HandleFileDTO) bool {
	pathsFile := strings.Split(dto.FilePath, "/")
	return len(pathsFile) > 4 && pathsFile[4] != ""
}

func (svc ParseFilesService) moveFileToDone(ctx context.Context, dto gross_measures.HandleFileDTO, err error) {
	stg, errStorage := svc.storage(ctx)
	if errStorage != nil {
		return
	}
	paths := strings.Split(dto.FilePath, "/")
	paths[0] = svc.storageSuccess
	if err != nil {
		paths[0] = svc.storageFail
	}
	dst := filepath.Join(paths...)
	err = stg.Copy(ctx, dto.FilePath, dst)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = stg.Delete(ctx, dto.FilePath)

	fmt.Println(err)
}

func (svc ParseFilesService) getHandler(ctx context.Context, dto gross_measures.HandleFileDTO) parsers.ParserCurveI {
	for _, h := range svc.handlers {
		if h.IsFile(ctx, dto) {
			return h
		}
	}

	return nil
}

func (svc ParseFilesService) getHandlerClose(ctx context.Context, dto gross_measures.HandleFileDTO) parsers.ParserCloseI {
	for _, h := range svc.handlersClose {
		if h.IsFile(ctx, dto) {
			return h
		}
	}

	return nil
}
