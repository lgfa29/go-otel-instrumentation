package main

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func startServer(bind string) error {
	// Initialize HTTP server.
	fmt.Printf("User registration server is running at %s\n", bind)

	wrappedHandler := otelhttp.NewHandler(http.HandlerFunc(signUp), "/")
	http.Handle("/", wrappedHandler)
	http.ListenAndServe(bind, nil)

	return nil
}

func signUp(w http.ResponseWriter, req *http.Request) {
	// Create span and ensure it ended at the end of the operation
	ctx := req.Context()

	operationName := "signing up"
	_, span := otel.Tracer("server").Start(ctx, operationName)
	defer span.End()

	// setting span "name" attribute
	name := req.FormValue("name")
	if name != "" {
		fmt.Println("name:", name)
		span.SetAttributes(attribute.String("name", name))
	} else {
		// setting span as error
		span.SetStatus(codes.Error, "missing name query parameter")
	}

	// setting span event
	span.AddEvent(fmt.Sprint(req.Header))

	// printing dummy registration
	w.Write([]byte("User " + name + " signed up"))
}
