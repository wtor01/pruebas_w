package log

import (
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

const (
	PROCESS_DURATION = "process_duration"
)

func New(service string) *otelzap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	c, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	logger := otelzap.New(c.With(
		zapcore.Field{
			Key:       "service",
			Type:      zapcore.StringType,
			Integer:   0,
			String:    service,
			Interface: nil,
		}))

	sugar := logger.Sugar()

	return sugar
}

func FieldDurationProcess() func() zap.Field {
	durationProcess := utils.DurationProcess()

	return func() zap.Field {
		return zap.Duration(PROCESS_DURATION, durationProcess())
	}

}

func NewDeprecated(service string) *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	logger = logger.With(zapcore.Field{
		Key:       "service",
		Type:      zapcore.StringType,
		Integer:   0,
		String:    service,
		Interface: nil,
	})

	sugar := logger.Sugar()

	return sugar
}
