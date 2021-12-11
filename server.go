package main

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
			semconv.ServiceNameKey.String("registration-server")))
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

	wrappedHandler := otelhttp.NewHandler(http.HandlerFunc(signUp), "/")
	http.Handle("/", wrappedHandler)
	fmt.Println("User registration server is running")
	http.ListenAndServe(":9000", nil)
}

func signUp(w http.ResponseWriter, req *http.Request) {

	// Create span and ensure it ended at the end of the operation
	ctx := req.Context()

	operationName := "signing up"
	_, span := otel.Tracer("server").Start(ctx, operationName)
	defer span.End()

	// before tracing
	name := "john"
	fmt.Print("name:", name)

	// after tracing
	name = "john"
	fmt.Println("name:", name)
	span.SetAttributes(attribute.String("name", name))

	// setting span as error
	span.SetStatus(codes.Error, "fatal")

	// setting span event
	span.AddEvent(fmt.Sprint(req.Header))

	// printing dummy registration
	w.Write([]byte("User " + name + " signed up"))

}
