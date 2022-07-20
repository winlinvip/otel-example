package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"os"
)

var tracer trace.Tracer
var provider *sdktrace.TracerProvider

func init() {
	////////////////////////////////////////////////
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		panic(err)
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("myService"),
		),
	)
	if err != nil {
		panic(err)
	}

	provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	fmt.Println("init ok")
}

func main() {
	// Create global tracer.
	tracer = otel.Tracer("app")
	defer provider.Shutdown(context.Background())

	// Create baggage.
	bag, err := baggage.Parse("myUserId=user0")
	if err != nil {
		panic(err)
	}

	// Create span with baggage.
	ctx := baggage.ContextWithBaggage(context.Background(), bag)
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	// Request HTTP server, with baggage.
	fmt.Println("wait to request...")

	foo(ctx)

	fmt.Println("otel done")
}

func foo(ctx context.Context) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8096", nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("body is", string(body))
}
