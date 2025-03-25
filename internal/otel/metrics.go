package otel

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func SetupOtelMeterProvider(
	ctx context.Context,
) (func(context.Context) error, error) {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint("otel-collector:4318"),
	}

	exporter, err := otlpmetrichttp.New(context.Background(), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "createing new metric exporter failed")
	}
	shutdownFunc := exporter.Shutdown

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new metrics resource failed: %w", err)
	}

	meterProvider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter, metricsdk.WithInterval(1))),
	)
	otel.SetMeterProvider(meterProvider)
	return shutdownFunc, nil
}

func GetMeter() metric.Meter {
	return otel.GetMeterProvider().Meter("")
}
