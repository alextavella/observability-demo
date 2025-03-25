package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func SetupOtelTracerProvider(
	ctx context.Context,
) (func(context.Context) error, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("otel-collector:4318"),
	}

	client := otlptracehttp.NewClient(opts...)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating new tracing exporter failed: %w", err)
	}
	shutdownFunc := exporter.Shutdown

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new tracing resource failed: %w", err)
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithResource(res),
		tracesdk.WithSampler(
			tracesdk.AlwaysSample(),
		),
		tracesdk.WithSpanProcessor(
			tracesdk.NewBatchSpanProcessor(exporter),
		),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return shutdownFunc, nil
}

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("")
}
