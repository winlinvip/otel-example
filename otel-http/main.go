package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"io"
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
	fmt.Println("init ok")
}

func main() {
	// Create global tracer.
	tracer = otel.Tracer("app")

	// See https://opentelemetry.io/docs/instrumentation/go/libraries/
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provider.ForceFlush(context.Background())
		io.WriteString(w, "Hello")
	}), "hello"))

	fmt.Println("Please test by http://localhost:8094")
	http.ListenAndServe(":8094", nil)
}
