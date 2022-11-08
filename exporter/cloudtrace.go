package exporter

import (
	"context"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type traceExporter struct {
	tp *sdktrace.TracerProvider
}

func Init(ctx context.Context, project, serviceNameKey string) (*traceExporter, error) {

	exporter, err := texporter.New(texporter.WithProjectID(project))
	if err != nil {
		return &traceExporter{}, err
	}

	res, err := resource.New(
		ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceNameKey)),
	)
	if err != nil {
		return &traceExporter{}, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	return &traceExporter{
		tp: provider,
	}, nil
}

func (te *traceExporter) Flush(ctx context.Context) {
	te.tp.ForceFlush(ctx)
}
