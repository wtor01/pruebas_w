package telemetry

import (
	"context"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
)

var tracer = otel.Tracer("app_or_package_name")

func tracerProviderJaeger(app string, url string) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(app),
		)),
	)
	return tp, nil
}

func tracerProvideGCP(app string) (*sdktrace.TracerProvider, error) {
	projectID := os.Getenv("PROJECT_ID")

	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.New: %v", err)
	}
	res, err := resource.New(context.Background(),
		// Use the GCP resource detector to detect information about the GCP platform
		resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application
		resource.WithAttributes(
			semconv.ServiceNameKey.String(app),
		),
	)
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}

	// Create trace provider with the exporter.
	//
	// By default it uses AlwaysSample() which samples all traces.
	// In a production environment or high QPS setup please use
	// probabilistic sampling.
	// Example:
	//   tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.0001)), ...)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

func NewTracerProvider(app string) *sdktrace.TracerProvider {
	var tp *sdktrace.TracerProvider
	var err error
	if os.Getenv("DEV") == "1" {
		tp, err = tracerProviderJaeger(app, "http://localhost:14268/api/traces")
	} else {
		tp, err = tracerProvideGCP(app)
	}

	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}
	otel.SetTracerProvider(tp)

	return tp
}

func SetTracer(name string) trace.Tracer {

	tracer = otel.GetTracerProvider().Tracer(name)

	return tracer
}

func GetTracer() trace.Tracer {
	return tracer
}
