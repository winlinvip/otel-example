package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var tracer trace.Tracer

func main() {
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

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	defer provider.Shutdown(context.Background())

	otel.SetTracerProvider(provider)
	fmt.Println("init ok")

	////////////////////////////////////////////////
	// Create global tracer.
	tracer = otel.Tracer("app")
	ctx, span := tracer.Start(context.Background(), "main", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	foo(ctx, "100")
	foo(ctx, "101")
	fmt.Println("done")
}

func foo(ctx context.Context, id string) {
	_, span := tracer.Start(ctx, "foo")
	defer span.End()

	span.AddEvent("log", trace.WithAttributes(attribute.KeyValue{Key: "id", Value: attribute.StringValue(id)}))
}
