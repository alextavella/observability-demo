package otel

import (
	"context"

	"github.com/pkg/errors"
)

func SetupOtel(
	ctx context.Context,
	local bool,
) (func(context.Context) error, error) {
	tracesShutdown, err := SetupOtelTracerProvider(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error setting up OTLP traces exporter")
	}
	metricsShutdown, err := SetupOtelMeterProvider(ctx)
	if err != nil {
		tracesShutdown(ctx)
		return nil, errors.Wrap(err, "error setting up OTLP metrics exporter")
	}

	return func(ctx context.Context) error {
		err := tracesShutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "error shutting down OTLP traces exporter")
		}
		err = metricsShutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "error shutting down OTLP metrics exporter")
		}
		return nil
	}, nil
}
