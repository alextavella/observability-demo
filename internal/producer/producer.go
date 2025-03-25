package producer

import (
	"context"
	"fmt"
	"log"
	"observability_demo/internal/otel"
	"time"

	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

	// Start the producer
	r := gin.Default()

	r.Use(otelgin.Middleware("producer"))

	r.GET("/", handleRoot)
	r.Run(":8082")
}

func handleRoot(c *gin.Context) {
	m := otel.GetMeter()
	counter, err := m.Int64Counter("producer.requests")
	if err != nil {
		fmt.Printf("error creating counter: %v\n", err)
	}

	time.Sleep(200 * time.Millisecond)

	v := c.Query("v")
	fmt.Println("Received request with value:", v)

	segmentioProducer := &kafka.Writer{
		Addr:  kafka.TCP("kafka:9092"),
		Topic: "opentel",
	}

	writer, err := otelkafkakonsumer.NewWriter(segmentioProducer,
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
			},
		))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer writer.Close()

	message := kafka.Message{
		Key:   []byte("v"),
		Value: []byte(v),
	}

	err = writer.WriteMessage(c.Request.Context(), message)

	if err == nil {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.Bool("success", true)))
		fmt.Println("Success!")
	} else {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.Bool("success", false)))
		fmt.Println("Failure!")
	}

	time.Sleep(100 * time.Millisecond)
}
