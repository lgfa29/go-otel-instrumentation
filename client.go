package main

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func main() {
	// START: Initializing tracing engine

	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// Service name to be use by observability tool
			semconv.ServiceNameKey.String("registration-client")))
	// Checking for errors
	if err != nil {
		fmt.Printf("Error adding %v to the tracer engine: %v", "applicationName", err)
	}

	collectorAddr := "127.0.0.1:1111"
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorAddr),
	)
	// Checking for errors
	if err != nil {
		fmt.Printf("Error initializing the tracer exporter: %v", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// END: initializing tracing engine

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport),
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9000", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("User registration request sent successfully")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}
