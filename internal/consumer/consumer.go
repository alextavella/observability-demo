package consumer

import (
	"context"
	"fmt"
	"observability_demo/internal/otel"
	"time"

	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Start() {
	ctx := context.Background()
	shutdown, err := otel.SetupOtel(ctx, false)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	reader, _ := otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{"kafka:9092"},
			GroupTopics: []string{"opentel"},
			GroupID:     "opentel-cg",
		}),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
			},
		),
	)

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("error reading message: %v\n", err)
			continue
		}
		processMessage(reader, message)
	}
}

func processMessage(reader *otelkafkakonsumer.Reader, message *kafka.Message) {
	m := otel.GetMeter()
	counter, err := m.Int64Counter("consumer.messages")
	if err != nil {
		fmt.Printf("error creating counter: %v\n", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Process the message
	v := string(message.Value)

	fmt.Printf("Received message: %s\n", v)

	if v != "" {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.Bool("success", true)))
		fmt.Println("Success!")
	} else {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.Bool("success", false)))
		fmt.Println("Failure!")
	}

	time.Sleep(100 * time.Millisecond)
}
