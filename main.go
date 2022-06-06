package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/detectors/nomad"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type OtelConfig struct {
	Address string
	Proto   string
}

func setupOtelTracing(serviceName string, config OtelConfig) (func(context.Context) error, error) {
	ctx := context.Background()

	serviceRes, err := resource.New(ctx,
		resource.WithAttributes(
			// Service name to be use by observability tool.
			semconv.ServiceNameKey.String(serviceName)))
	if err != nil {
		return nil, fmt.Errorf("Error adding %v to the tracer engine: %v", "applicationName", err)
	}

	nomadRes, err := nomad.NewResourceDetector().Detect(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error adding Nomad to the tracer engine: %v", err)
	}

	res, err := resource.Merge(serviceRes, nomadRes)
	if err != nil {
		return nil, fmt.Errorf("Error merging tracer engines: %v", err)
	}

	var traceExporter *otlptrace.Exporter
	switch config.Proto {
	case "http":
		traceExporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(config.Address),
		)
	case "grpc":
		traceExporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Address),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("Error initializing the tracer exporter: %v", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	fmt.Printf("Connecting to the OTel Collector on %s via %s\n", config.Address, config.Proto)

	return tp.Shutdown, nil
}

func getEnvDefault(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return v
}

func main() {
	ctx := context.Background()

	otelAddr := getEnvDefault("OTEL_COLLECTOR_ADDR", "127.0.0.1:1111")
	otelProto := getEnvDefault("OTEL_COLLECTOR_PROTO", "http")
	otelConfig := OtelConfig{
		Address: otelAddr,
		Proto:   otelProto,
	}

	if len(os.Args) < 2 {
		fmt.Printf("missing required argument mode, either client or server must be provided.\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		shutdown, err := setupOtelTracing("registration-server", otelConfig)
		if err != nil {
			fmt.Printf("failed to start tracing: %v\n", err)
			os.Exit(1)
		}
		defer shutdown(ctx)

		bind := getEnvDefault("BIND_ADDR", ":9000")
		err = startServer(bind)
		if err != nil {
			fmt.Printf("failed to start server: %v\n", err)
			os.Exit(1)
		}
	case "client":
		clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
		name := clientCmd.String("name", "Kathryn Janeway", "")
		clientCmd.Parse(os.Args[2:])

		shutdown, err := setupOtelTracing("registration-client", otelConfig)
		if err != nil {
			fmt.Printf("failed to start tracing: %v\n", err)
			os.Exit(1)
		}
		defer shutdown(ctx)

		url := getEnvDefault("SERVER_URL", "http://127.0.0.1:9000")
		err = makeRequest(ctx, url, *name)
		if err != nil {
			fmt.Printf("failed to run client: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("invalid mode %s: only client or server\n", os.Args[1])
		os.Exit(1)
	}
}
