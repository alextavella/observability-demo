package backend

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"observability_demo/internal/otel"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Start() {
	ctx := context.Background()
	shutdown, err := otel.SetupOtel(ctx, false)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	// Start the backend
	r := gin.Default()

	r.Use(otelgin.Middleware("backend"))

	r.GET("/", handleRoot)
	r.Run(":8081")

}

func handleRoot(c *gin.Context) {
	m := otel.GetMeter()
	counter, err := m.Int64Counter("backend.requests")
	if err != nil {
		fmt.Printf("error creating counter: %v\n", err)
	}

	t := otel.GetTracer()
	time.Sleep(100 * time.Millisecond)

	_, span1 := t.Start(c.Request.Context(), "db query")
	time.Sleep(100 * time.Millisecond)
	span1.End()

	v := c.Query("v")
	fmt.Println("Received request with value:", v)

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
			}),
		),
	}

	req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, fmt.Sprintf("http://producer:8082?v=%s", v), nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error sending request to producer: %v\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(res.StatusCode)

	if res.StatusCode == http.StatusOK {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.Bool("success", true)))
		fmt.Println("Success!")
	} else {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.Bool("success", false)))
		fmt.Println("Failure!")
	}

	_, span2 := t.Start(c.Request.Context(), "db save")
	time.Sleep(100 * time.Millisecond)
	span2.End()

	time.Sleep(100 * time.Millisecond)
}
