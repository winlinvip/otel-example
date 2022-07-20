package main

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

func main() {
	// Parse span from request.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// See https://opentelemetry.io/docs/instrumentation/go/libraries/
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		fmt.Println("trace is", span.SpanContext().TraceID().String())

		bag := baggage.FromContext(r.Context())
		fmt.Println("bag is", bag)

		fmt.Println("Client header", r.Header)

		io.WriteString(w, "Hello")
	}), "hello"))

	fmt.Println("Please test by http://localhost:8096")
	http.ListenAndServe(":8096", nil)
}
