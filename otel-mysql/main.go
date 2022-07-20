package main

import (
	"context"
	"fmt"
	"github.com/signalfx/splunk-otel-go/instrumentation/database/sql/splunksql"
	_ "github.com/signalfx/splunk-otel-go/instrumentation/github.com/go-sql-driver/mysql/splunkmysql"
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

	////////////////////////////////////////////////////////////////////////
	// Connect to db, see https://github.com/signalfx/splunk-otel-go/blob/main/instrumentation/github.com/go-sql-driver/mysql/splunkmysql/example_test.go
	db, err := splunksql.Open("mysql", "root:12345678@tcp(127.0.0.1:13306)/mysql")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////
	// Flush otel manually.
	provider.ForceFlush(context.Background())

	// See https://github.com/signalfx/splunk-otel-go/blob/main/instrumentation/github.com/go-sql-driver/mysql/splunkmysql/example_test.go
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var now string
		if err := db.QueryRowContext(r.Context(), "select CURRENT_TIMESTAMP").Scan(&now); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		provider.ForceFlush(context.Background())
		io.WriteString(w, "Hello "+now)
	}), "hello"))

	fmt.Println("Please test by http://localhost:8095")
	http.ListenAndServe(":8095", nil)
}
